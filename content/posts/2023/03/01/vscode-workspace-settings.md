---
title: "ใช้ VSCode settings เดียวกันทั้งโปรเจคด้วย Workspace settings"
date: 2023-03-01T07:32:46+07:00
draft: false
---

ถ้าทีมใช้ VSCode เหมือนๆกันอยู่แล้ว การใช้ Workspace settings จะช่วยให้ทั้งทีมที่ดูแลโปรเจคเดียวกันอยู่ใช้ค่า settings ที่เหมือนกันทั้งโปรเจคได้

<!--more-->

วิธีสร้าง workspace settings ก็ง่ายๆ ให้เราเพิ่ม directory `.vscode` และสร้างไฟล์ config `settings.json` ในนั้น ตัวอย่างเช่นเรามีโปรเจค `todoapp` เราก็สร้างไฟล์ไว้ใน `.vscode/settigns.json` ไว้ใน `todoapp` แบบนี้

```
├── todoapp
    ├── .vscode
    │   └──settings.json
    ├── app
```

ถ้าเราเข้าไปใน settings ของ VSCode เราก็จะเห็น tab Workspace ซึ่งถ้าเราเซตค่าจาก UI ก็จะ save ลงมาที่ไฟล์ `.vscode/settings.json` นั่นเอง

ตัวอย่างที่ผมใช้กับทีมก็เช่นพวก settings การ format ในแต่ละภาษาที่ใช้กันในโปรเจค และ Playwright ENV

```json
{
  "editor.codeActionsOnSave": {
    "source.fixAll": true
  },
  "editor.formatOnSave": true,
  "[typescript]": { "editor.defaultFormatter": "esbenp.prettier-vscode" },
  "[typescriptreact]": { "editor.defaultFormatter": "esbenp.prettier-vscode" },
  "playwright.reuseBrowser": true,
  "playwright.env": {},
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  }
}
```

ลองเอาไปใช้กันดูครับ ช่วยให้ทีมทำงานด้วยกันง่ายขึ้นเยอะเลย ส่วนใครใช้ IDE/Editor อื่นก็คงจะมี settings ในระดับ project/workspace แบบนี้เช่นกัน ก็ควรปรับให้ตรงกันด้วยเหมือนกันครับ

ref:

- https://code.visualstudio.com/docs/getstarted/settings#_workspace-settings
- https://code.visualstudio.com/docs/editor/workspaces
