---
title: "[Go] สร้าง test context type เพื่อให้โค้ดของ test เรียบง่ายขึ้น"
date: 2023-03-02T07:09:48+07:00
draft: false
---

ในโค้ดของเทส หลายครั้งมีการทำงานในส่วนของการเตรียมข้อมูลก่อนรันเทส และการเคลียร์ข้อมูลหลังจากเทสรันเสร็จที่ซ้ำซ้อนกัน เราสามารถแยกออกมาเพื่อให้โค้ดเทสเรียบง่ายขึ้นได้ วันนี้จะใช้เทคนิคสร้าง type ใหม่เรียกว่า test context เพื่อยุบสิ่งที่ซ้ำซ้อนมาเป็น method ของ type นี้แทน

<!--more-->

ตัวอย่างโค้ด เป็นโค้ดที่เขียนต่อกับ Firebase Auth ผ่าน [Firebase Admin Go SDK v4](https://pkg.go.dev/firebase.google.com/go/v4) ซึ่งเราก็มี method ที่จะเทสเช่น get user by UID, activate account by UID, deactivate account by UID

ทีนี้ตอนเขียนเทส ก็ต้องมีการสร้าง firebase auth client, ต้อง create user เพื่อเทส ซึ่งก็จะมีของซ้ำๆกันอยู่ แบบนี้

```go
func TestGetUserByUID(t *testing.T) {
	ctx := context.Background()
	fbApp, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: "test",
	})
	require.NoError(t, err)

	fbAuth, err := New(ctx, fbApp)
	require.NoError(t, err)
	client := fbAuth.client

	userRecord, err := createUser(ctx, client, "admin@example.com", "password1234", "admin")
	require.NoError(t, err)
	defer func() {
		err := client.DeleteUser(ctx, userRecord.UID)
		require.NoError(t, err)
	}()

	user, err := fbAuth.GetUserByUID(ctx, userRecord.UID)
	require.NoError(t, err)

	expectedUser := &User{
		UID:      userRecord.UID,
		Email:    "admin@example.com",
		Role:     "admin",
		Activate: true,
	}
	require.Equal(t, expectedUser, user)
}

func TestActivateUser(t *testing.T) {
	ctx := context.Background()
	fbApp, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: "test",
	})
	require.NoError(t, err)

	fbAuth, err := New(ctx, fbApp)
	require.NoError(t, err)
	client := fbAuth.client

	userRecord, err := createUser(ctx, client, "planner@example.com", "password1234", "planner")
	require.NoError(t, err)
	defer func() {
		err := client.DeleteUser(ctx, userRecord.UID)
		require.NoError(t, err)
	}()

	userRecord, err = client.UpdateUser(ctx, userRecord.UID, (&auth.UserToUpdate{}).Disabled(true))
	require.NoError(t, err)
	require.True(t, userRecord.Disabled)

	user, err := fbAuth.ActivateUser(ctx, userRecord.UID)
	require.NoError(t, err)

	expectedUser := &User{
		UID:      userRecord.UID,
		Email:    "planner@example.com",
		Role:     "planner",
		Activate: true,
	}
	require.Equal(t, expectedUser, user)
}

func TestDeactivateUser(t *testing.T) {
	ctx := context.Background()
	fbApp, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: "test",
	})
	require.NoError(t, err)

	fbAuth, err := New(ctx, fbApp)
	require.NoError(t, err)
	client := fbAuth.client

	userRecord, err := createUser(ctx, client, "planner@example.com", "password1234", "planner")
	require.NoError(t, err)
	defer func() {
		err := client.DeleteUser(ctx, userRecord.UID)
		require.NoError(t, err)
	}()
	require.False(t, userRecord.Disabled)

	user, err := fbAuth.DeactivateUser(ctx, userRecord.UID)
	require.NoError(t, err)

	expectedUser := &User{
		UID:      userRecord.UID,
		Email:    "planner@example.com",
		Role:     "planner",
		Activate: false,
	}
	require.Equal(t, expectedUser, user)
}

func createUser(ctx context.Context, client *auth.Client, email, password, role string) (*auth.UserRecord, error) {
	createUserParam := &auth.UserToCreate{}
	createUserParam.Email(email)
	createUserParam.Password(password)
	userRecord, err := client.CreateUser(ctx, createUserParam)
	if err != nil {
		return nil, err
	}

	err = client.SetCustomUserClaims(ctx, userRecord.UID, map[string]any{"role": role})
	if err != nil {
		return nil, err
	}

	return client.GetUser(ctx, userRecord.UID)
}
```

จะเห็นว่าแต่ละเคสเรามีขั้นตอนการเตรียม firebase client ซ้ำซ้อนกัน แม้ว่าตัวอย่างโค้ดที่จะยุบ createUser ออกเป็น helper function แล้วก็ตาม ก็ยังดูรกๆในขั้นตอนเตรียมการอยู่

ทีนี้เราจะ refactor กันใหม่โดยสร้าง type testContext struct ขึ้นมาเพื่อเก็บสิ่งที่จะแชร์ร่วมกันในแต่ละเทสแล้วก็ method helper ที่จะช่วยยุบให้โค้ดของเทสเรียบง่ายยิ่งขึ้น

```go
type testContext struct {
	ctx        context.Context
	fba        *firebaseAuth
	authClient *auth.Client
	t          *testing.T
}

func newTestContext(t *testing.T) *testContext {
	ctx := context.Background()
	fbApp, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: "test",
	})
	require.NoError(t, err)

	fbAuth, err := New(ctx, fbApp)
	require.NoError(t, err)

	return &testContext{
		ctx:        ctx,
		authClient: fbAuth.client,
		t:          t,
		fba:        fbAuth,
	}
}

func (testCtx *testContext) createUser(email, password, role string) *auth.UserRecord {
	createUserParam := &auth.UserToCreate{}
	createUserParam.Email(email)
	createUserParam.Password(password)
	userRecord, err := testCtx.authClient.CreateUser(testCtx.ctx, createUserParam)
	require.NoError(testCtx.t, err)

	err = testCtx.authClient.SetCustomUserClaims(testCtx.ctx, userRecord.UID, map[string]any{"role": role})
	require.NoError(testCtx.t, err)

	userRecord, err = testCtx.authClient.GetUser(testCtx.ctx, userRecord.UID)
	require.NoError(testCtx.t, err)

	testCtx.t.Cleanup(func() {
		err := testCtx.authClient.DeleteUser(testCtx.ctx, userRecord.UID)
		require.NoError(testCtx.t, err)
	})
	return userRecord
}

func (testCtx *testContext) createDisabledUser(email, password, role string) *auth.UserRecord {
	user := testCtx.createUser(email, password, role)
	userRecord, err := testCtx.authClient.UpdateUser(testCtx.ctx, user.UID, (&auth.UserToUpdate{}).Disabled(true))
	require.NoError(testCtx.t, err)
	require.True(testCtx.t, userRecord.Disabled)
	return userRecord
}
```

จะเห็นว่าเราสร้าง tyep testContext struct แล้วจับ ค่าต่างๆที่มีร่วมกันหลายๆเทสมาเป็น field ใน type นี้ จากนั้นสร้าง newTestContext ที่รวบรวมกัน setup ค่าต่างๆที่ซ้ำซ้อนกันเอาไว้ ย้ายการ createUser ทั้งแบบ active, inactive มาเป็น method ของ testContext ซึ่งก็ทำให้ method ไม่ต้องรับ parameter เยอะเกินเพราะเราเก็บค่าจำเป็นใน field ของ testContext struct แล้ว

ทีนี้มาดูโค้ดหลังจากเอา testContext ไปใช้งานกัน

```go
func TestGetUserByUID(t *testing.T) {
	testCtx := newTestContext(t)
	userRecord := testCtx.createUser("admin@example.com", "password1234", "admin")

	user, err := testCtx.fba.GetUserByUID(testCtx.ctx, userRecord.UID)
	require.NoError(t, err)

	expectedUser := &User{
		UID:      userRecord.UID,
		Email:    "admin@example.com",
		Role:     "admin",
		Activate: true,
	}
	require.Equal(t, expectedUser, user)
}

func TestActivateUser(t *testing.T) {
	testCtx := newTestContext(t)
	userRecord := testCtx.createDisabledUser("planner@example.com", "password1234", "planner")

	user, err := testCtx.fba.ActivateUser(testCtx.ctx, userRecord.UID)
	require.NoError(t, err)

	expectedUser := &User{
		UID:      userRecord.UID,
		Email:    "planner@example.com",
		Role:     "planner",
		Activate: true,
	}
	require.Equal(t, expectedUser, user)
}

func TestDeactivateUser(t *testing.T) {
	testCtx := newTestContext(t)
	userRecord := testCtx.createUser("planner@example.com", "password1234", "planner")

	user, err := testCtx.fba.DeactivateUser(testCtx.ctx, userRecord.UID)
	require.NoError(t, err)

	expectedUser := &User{
		UID:      userRecord.UID,
		Email:    "planner@example.com",
		Role:     "planner",
		Activate: false,
	}
	require.Equal(t, expectedUser, user)
}
```

โค้ดเรียบง่ายขึ้นมาก เพราะเราแค่ newTestContext ขึ้นมา แล้วเรียก method helper อย่าง createUser เสร็จแล้วเรียก method ของ type ที่เราต้องการเทส แล้วก็ทำการ verify ผลลัพธ์

ลองไปค้นดูเล่นๆ ก็มีโปรเจคอื่นๆใน pattern แบบนี้อยู่เหมือนกันเช่น gvisor มี type ที่ช่วยเก็บ state ในการเทสเช่นกันแบบนี้ (https://cs.opensource.google/gvisor/gvisor/+/master:pkg/tcpip/tests/integration/istio_test.go;l=41;drc=1338761211656f1a1b3cb358fcf778ae02a0a4ec)

```go
// testContext encapsulates the state required to run tests that simulate
// an istio like environment.
type testContext struct {
```

โปรเจ็ค teleport (https://github.com/gravitational/teleport/blob/e28d969fdcc9a7c9af9f7cb4ef2281b390fe3a95/lib/kube/proxy/exec_test.go#L59)

ก็มี TestContext เช่นกัน

```go
	testCtx := SetupTestContext(
		context.Background(),
		t,
		TestConfig{
			Clusters: []KubeClusterConfig{{Name: kubeCluster, APIEndpoint: kubeMock.URL}},
		},
	)
```

ใครที่มีโค้ดเทสที่ต้องเซตอัพค่าอะไรหลายอย่างและใช้ซ้ำๆกันหลายๆเทสเคส ลองเอาท่านี้ไปประยุกต์ใช้กันดูครับ
