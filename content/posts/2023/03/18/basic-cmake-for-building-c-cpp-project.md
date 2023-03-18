---
title: "[CMake] config Cmake เพื่อ build C/C++ เบื้องต้น"
date: 2023-03-18T17:00:00+07:00
draft: false
---

[CMake](https://cmake.org) เป็น build tool นึงเพื่อช่วยให้เราเซตโปรเจ็ค C/C++ ให้ build binary ออกมาได้ง่ายๆ เพราะ C/C++ เวลาจะ build มีทั้งต้องกำหนด include ตอน compile ทั้งต้องกำหนด library อื่นๆที่เกี่ยวข้องก่อน linking ดังนั้นใช้ CMake ก็จะช่วยจัดการพวกนี้ได้ง่ายขึ้น

<!--more-->

สิ่งที่เราต้องทำคือสร้างไฟล์ `CMakeLists.txt` ซึ่งเป็นไฟล์ config ของ CMake สำหรับการ build C/C++ เบื้องต้น ที่มีใช้ Boost library ด้วย เป็นแบบนี้

```txt
# กำหนด CMake minimum version ไม่ต่ำกว่า 3.18
cmake_minimum_required(VERSION 3.18)

# กำหนดชื่อโปรเจค boost-thread-pool-example และภาษาที่ใช้คือ C และ CXX (C++)
project(boost-thread-pool-example C CXX)

# เช็คว่าถ้า env CMAKE_BUILD_TYPE ไม่ใช่ Debug ก็จะกำหนด CMAKE_BUILD_TYPE เป็น Release
if(NOT CMAKE_BUILD_TYPE MATCHES Debug)
  set(CMAKE_BUILD_TYPE Release)
endif()

# กำหนด version ว่าจะใช้ CPP version 17
set(CMAKE_CXX_STANDARD 17)
# เปิดใช้งาน flag CMAKE_CXX_STANDARD_REQUIRED
set(CMAKE_CXX_STANDARD_REQUIRED ON)
# ปิดการใช้งาน CPP Extensions
set(CMAKE_CXX_EXTENSIONS OFF)

# เช็คว่า compile บน 64 bit หรือ 32 bit
if(CMAKE_SIZEOF_VOID_P EQUAL 8)
  message(STATUS "Building on a 64 bit system")
else()
  message(STATUS "Building on a 32 bit system")
endif()

# จะใช้คำสั่ง find_package เพื่อระบุ dependencies เช่นกรณีนี้คือ Boost library
# ซึ่ง CMake จะ include module FindBoost.cmake เขามาให้เราด้วย
# ในนั้นก็จะหา include path, lib path แล้วกำหนดค่าลงตัวแปรที่ขึ้นต้นด้วย Boost_ ให้เราเช่น
# Boost_CXXFLAGS
find_package(Boost REQUIRED)

# เพิ่ม include directory ของ Boost เข้ามาใน project
include_directories(${Boost_INCLUDE_DIRS})

# ระบุว่า เราจะ build executable program เราจาก main.cpp และชื่อไฟล์โปรแกรมคือ boost-thread-pool-example
add_executable(boost-thread-pool-example main.cpp)

# สุดท้ายตอน link library ระบุด้วยว่าตอน link ให้เอา library ของ Boost มา link ด้วย
target_link_libraries(boost-thread-pool-example ${Boost_LIBRARIES})
```

ทีนี้ในโปรเจ็คเราตอนจะ build ก็ให้สร้าง build directory ขึ้นมาก่อนเช่นในโปรเจ็คมีไฟล์แบบนี้

```txt
app
- main.cpp
- CMakeLists.txt
```

ให้เรา cd เข้าไปใน `app` แล้ว `mkdir build` แล้ว cd เข้าไปใน `build` แล้วสั่ง `cmake ..` แบบนี้

```txt
cd app
makdir build
cd build
cmake ..
```

เสร็จแล้วตอนจะ build ก็ให้สั่ง

```
cmake --build .
```

ก็ได้แล้ว

ถ้าใครใช้ VSCode CMake extension ก็สั่ง build จาก VSCode ได้เลยอีกด้วย
