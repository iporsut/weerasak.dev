---
title: สรุปจากเรื่อง Livable Code ของ Sarah Mei จากงาน RailsConf 2018
date: 2023-06-18T19:27:00+07:00
draft: false
---

เคยเขียนสรุปสิ่งที่เข้าใจและได้จากการฟังคลิปเรื่อง Livable Code เอาไว้แล้วรอบนึงที่เพจ [DevDose](https://www.facebook.com/devdoseth/posts/pfbid0YBoKMVdB5w1xRHvbgRH5CSvAZzoPi5L3n1YAk3CJ74eJFfk5anvn6tHoeEcg8jZql) ขอยกมาเก็บไว้ใน blog ด้วยจะได้อีกที เป็น session ที่ชอบมากๆหนึ่ง session

<!--more-->

อยากให้ทุกคนในวงการพัฒนาซอฟแวร์ได้ดูได้ฟัง Talk นี้จาก Sarah Mei มากๆ ผมสรุปจากที่ฟังเอาไว้คร่าวๆ อย่าลืมเข้าไปดูตัวเต็มกันใน ลิ้งนี้นะครับ Livable Code => https://www.youtube.com/watch?v=lI77oMKr5EY

> Build trust between team members<br/>
> and strengthen those connections.<br/>
> The code base will follow.<br/>

Sarah Mei เล่าให้เล่าฟังว่า Software และ การพัฒนา Software ปัจจุบันไม่สามารถอธิบายได้โดยใช้การเทียบเคียงกับ กระบวนการพัฒนาโปรดักอื่นๆหรืองานก่อสร้างที่แบ่งเป็นเฟส Architecture => Engineering => Construction อีกต่อไปแล้ว

เธอเสนอว่า เราควรมอง Software เนี่ยเป็น Livable Code หรือ Livable Product งานซอฟแวร์ส่วนใหญ่ในปัจจุบันเราไม่ได้มองเหมือนการสร้างตึกที่เริ่มทุกอย่างจากพื้นที่ดินว่างๆกันอีกต่อไปแล้ว เราพัฒนาซอฟแวร์กันเหมือนการออกแบบภายในมากกว่า ตัวตึก ตัวโครงสร้าง infrastructure ต่างๆมีไว้ให้แล้ว

เทียบได้กับปัจจุบันตัวซอฟแวร์เราสร้างโดยเราเลือก framework เป็นโครงสร้างแล้วเราตกแต่งต่อเติมจากนั้นอีกที หรืออย่างทำเกม เราก็มี Game Engine เป็นโครงสร้างหลักให้เรา และเธอมองว่าเราไม่ได้จัดโครงสร้างภายในเพื่อแค่ถ่ายรูปโชว์ เราจัดบ้านก็เพื่อให้คนอยู่อาศัยร่วมกัน หรือแม้แต่จัดเสร็จแล้ว คนที่อยู่ในบ้านเปลี่ยนไปเดี๋ยวสิ่งต่างๆ การจัดบ้านก็เปลี่ยนไป ตอนอยู่กันสองคนก็อีกแบบ ครอบครัวใหญ่ขึ้นมีลูกก็อีกแบบ

ซึ่งก็เหมือนกับ Code base ที่เราเขียนกัน จะเป็นยังไงก็ขึ้นอยู่กับคนที่ใช้ชีวิตอยู่ร่วมกันกับมัน มันไม่มีสูตรตายตัวว่าจัดบ้านจัดโต๊ะแบบไหนถึงดี บางคนชอบโต๊ะรกๆ แต่ตัวเองรู้ว่าจะของอะไรอยู่ตรงไหน บางคนชอบจัดของเป็นระเบียบ แต่คนอยู่ด้วยอาจจะไม่ชอบเพราะหาของยาก

เธอเสนอแนวทาง 4 ข้อว่าจะปรับปรุงของในบ้าน (โค้ด) ได้แก่

1. Don't Make it Worse
   อย่าทำให้มันแย่กว่าเดิม

2. Improvement Over Consistency
   ค่อยๆ improve ไป

3. Inline Everything
   คล้ายกับข้อที่แล้ว คือการปรับปรุงไม่ต้องรอปรับใหญ่ ปรับไปเลยทีละนิดในจุดที่ทำได้

4. Talk More
   คุยกันเยอะๆ <br />
   4a. Don't ask permission but be upfront
   ไม่ต้องขออนุญาต ถ้าจะปรับปรุงอะไรให้มันดีขึ้น แต่ก็แจ้งล่วงหน้ากันด้วย <br />
   4b. Don't ask forgiveness
   ไม่ต้องขอให้ยกโทษอะไร ไม่ได้ทำอะไรผิด <br />
   4c. Do ask for advice But don't always take it
   ถามเพื่อขอคำปรึกษา แต่ไม่ต้องเอาทุกคำแนะนำมาปฎิบัติ <br />
   4d. Do work together You have to live here
   ทำงานด้วยกันปรับปรุงด้วยกัน เพราะเราต้องอยู่ด้วยกันที่บ้าน (อยู่กับโค้ดเบสเดียวกัน)

และบอกว่า สิ่งสำคัญสุดใน Software ปัจจุบันไม่ใช่แค่ Code ไม่ใช่แค่ People แต่คือทั้งสองอย่างที่รวมกันเป็น System
และเธอทิ้งท้ายใน talk นี้เอาไว้ว่า

> Build trust between team members<br/>
> and strengthen those connections.<br/>
> The code base will follow.<br/>