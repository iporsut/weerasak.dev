---
title: "[C++] ใช้ CPU ให้เต็มที่ด้วย Boost Thread Pool"
date: 2023-03-11T08:40:00+07:00
draft: false
---

เมื่อวานมีงานที่ต้อง process ข้อมูลบางอย่างด้วย C++ เลยพยายามใช้ CPU ให้เต็มที่ทุก cores เท่าที่มีให้มากสุด ซึ่งก็ได้มาท่านึงก็คือใช้ library Boost Thread Pool เข้าช่วย

<!--more-->

โจทย์สมมติว่า เรามีข้อมูล n ตัว เราต้องการโปรเซสข้อมูลของการจับคู่ของมูลเช่น

```txt
n = [1, 2, 3]
```

แล้วเรามีฟังก์ชัน process ข้อมูล ที่รับสองพารามิเตอร์แบบนี้

```cpp
void process(int a, int b)
{
        cout << a << "+" << b << " = " << a + b << "\n";
}
```

ตัวอย่างโค้ดแบบไม่ได้ใช้ Thread Pool ช่วยจะเป็นแบบนี้

```cpp
#include <iostream>
#include <vector>

using namespace std;

void process(int a, int b)
{
        cout << a << "+" << b << " = " << a + b << "\n";
}

int main()
{
        vector<int> vs = {1, 2, 3, 4, 5};

        for (const auto &a : vs)
        {
                for (const auto &b : vs)
                {
                        process(a, b);
                }
        }
        return EXIT_SUCCESS;
}
```

(ป.ล. ของจริงที่ใช้ไม่ใช่แค่ print เลขบวกกัน แต่เมื่อให้เข้าใจง่ายเฉยๆ ต้องจินตนาการว่า process เป็นงานที่คำนวณหนักๆหน่อย :D)

ทีนี้เราอยากให้ตอน process แทนทีจะทำทีละ 1 แบบ sequential เราอยากใช้ thread ช่วยเพื่อให้ใช้ cores CPU ได้เต็มที่ เราเอา Boost มาช่วยได้โดยโค้ดจะเป็นแบบนี้แทน

```cpp
#include <boost/asio.hpp>
#include <boost/thread/mutex.hpp>

#include <iostream>
#include <vector>
#include <thread>
#include <functional>

using namespace std;

void process(boost::mutex &io_mutex, int a, int b)
{
        boost::mutex::scoped_lock scoped_lock(io_mutex);
        cout << a << "+" << b << " = " << a + b << "\n";
}

int main()
{
        // สร้าง thread_pool object เก็บในตัวแปร pool
        // โดยกำหนดให้จำนวน thread เท่ากับ core ของ CPU ของเครื่องที่รัน
        // เช็คได้จาก function hardward_concurrency
        boost::asio::thread_pool pool(thread::hardware_concurrency());

        // สร้าง mutex เพื่อช่วยให้ตอน cout ไม่เกิด race condition ที่ทำให้ผลลัพธ์เพี้ยนๆ
        boost::mutex io_mutex;

        vector<int> vs = {1, 2, 3, 4, 5};

        for (const auto &a : vs)
        {
                for (const auto &b : vs)
                {
                        // ส่งฟังก์ชันไปรันใน thread pool ด้วยฟังก์ชัน post
                        // ซึ่งเราต้องใช้ std::bind ฟังก์ชันกับ parameter ของฟังก์ชันก่อนส่งให้ post ด้วย
                        // และส่ง reference io_mutex เข้าไปด้วยเพื่อใช้ scoped_lock กันแย่งกัน print output
                        boost::asio::post(pool, std::bind(process, std::ref(io_mutex), a, b));
                }
        }

        // เรียก pool.join() เพื่อรอให้ทุก process ทำงานจบก่อนแล้วค่อยจบโปรแกรม
        pool.join();

        return EXIT_SUCCESS;
}
```

จะเห็นว่าโค้ดนี้มีจุดน่าสนใจเช่นตรง `boost::asio::post(pool, std::bind(process, std::ref(io_mutex), a, b));` เราได้ใช้ฟีเจอร์ใหม่ๆของ C++ ที่เอาฟังก์ชันที่ประกาศไว้มาทำให้ pass เป็น value ของ post โดยใช้ `std::bind` ซึ่งส่งฟังก์ชันกับลิสต์ของพารามิเตอร์ไปให้ ซึ่งค่าไหนเป็น reference ก็ให้ใช้ `std::ref` ครอบก่อน

นอกจากนั้นเรายังได้ใช้ mutex กับ scoped_lock ของ boost::asio ช่วยในการจัดการ lock ตอนแสดงผลด้วย cout ด้วยเพราะถ้าไม่ทำจะเกิด race condition ตอนแสดงผล จนผลลัพธ์ผิดเพี้ยนไปได้

scoped_lock นั้นใช้วิธี RAII ซึ่งจะ unlock เองตอนจบ scope ทำให้เราไม่ต้องเรียก unlock เองอีกด้วย แต่สร้าง scoped_lock ก็พอแล้ว

ใน C++ ใหม่ๆรองรับ lambda express ด้วย ดังนั้นเดี๋ยวเราจะลองแปลงโค้ดใหม่มาใช้ lambda ดูแทนที่จะใช้ฟังก์ชัน std::bind process โค้ดก็จะได้แบบนี้

```cpp
#include <boost/asio.hpp>
#include <boost/thread/mutex.hpp>

#include <iostream>
#include <vector>
#include <thread>
#include <functional>

using namespace std;

void process(boost::mutex &io_mutex, int a, int b)
{
        boost::mutex::scoped_lock scoped_lock(io_mutex);
        cout << a << "+" << b << " = " << a + b << "\n";
}

int main()
{
        // สร้าง thread_pool object เก็บในตัวแปร pool
        // โดยกำหนดให้จำนวน thread เท่ากับ core ของ CPU ของเครื่องที่รัน
        // เช็คได้จาก function hardward_concurrency
        boost::asio::thread_pool pool(thread::hardware_concurrency());

        // สร้าง mutex เพื่อช่วยให้ตอน cout ไม่เกิด race condition ที่ทำให้ผลลัพธ์เพี้ยนๆ
        boost::mutex io_mutex;

        vector<int> vs = {1, 2, 3, 4, 5};

        for (const auto &a : vs)
        {
                for (const auto &b : vs)
                {

                        // ส่งฟังก์ชันไปรันใน thread pool ด้วยฟังก์ชัน post
                        // โดยส่ง lambda express ที่เรียก process ข้างใน
                        boost::asio::post(
                            pool,
                            // lambda express syntax
                            // แบบนี้คือ [&] หมายถึงทุกตัวแปรที่ capture จาก scope ด้านนอกใช้เป็น reference
                            // ทั้ง io_mutex, a, b โดยเราไม่ต้องลิสต์เองทีละอัน
                            [&]()
                            {
                                    process(io_mutex, a, b);
                            });
                        // boost::asio::post(pool, std::bind(process, std::ref(io_mutex), a, b));
                }
        }

        // เรียก pool.join() เพื่อรอให้ทุก process ทำงานจบก่อนแล้วค่อยจบโปรแกรม
        pool.join();

        return EXIT_SUCCESS;
}
```

จะเห็นว่าโค้ด C++ ยุคใหม่ๆ ก็ทำให้อ่านและเขียนง่ายขึ้นเยอะเลย ใครมีงานอะไรที่ต้องใช้ lib c, lib cpp ลองมาเขียน C++ ใช้งานดูก็ไม่ได้ดูยากอีกต่อไปแล้วนะ
