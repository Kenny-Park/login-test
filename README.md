 # login 구현
로그인


## 실행방법
* /etc 폴더에 있는 DDL.sql 스크립트를 mysql에서 실행한다.
* 레파지토리에 올라간 login_test 파일을 터미널에서 실행시킨다.(맥북기준)
* 브라우저에서 localhost:8080/login으로 실행한다.

## 사용기술

구분|기술
---|---|
언어|GO 1.16.3|
DB|mysql|
라이브러리|xorm, gin-gonic framework|

## 구현특이사항
* 토큰방식으로 구현

프론트엔드
* JQUERY 활용
* 세션,쿠키대신 만료시간이 존재하는 로컬스토리지를 구현하여 활용

백엔드
* restful api
* 정보 암호화 적용 (개인정보 AES, 비밀번호 SHA256)
* xorm 라이브러리 활용
* micro web framework(gin-gonic) 활용

API
URL|METHOD|PARAM|REQ|RESP|비고
---|---|---|---|---|---|
/usr/auth/phone/req|POST|hp=01000000000&userid=pkhloved| - |{"code": 200, "message":"", "trankey": "20220826160000SKLDSAJDS", "num": "111111"}|인증번호요청, num은 휴대폰인증을 할 수 없으므로 확인용도로 내려줌|
/usr/auth/phone/conf|POST|hp=01000000000&userid=pkhloved&trankey=20220826160000SKLDSAJDS&num=111111| - |{"code": 200, "message":""}|인증번호확인|
/usr/login|POST|loginid=pkhloved&pwd=dkagh| - |{"code": 200, "message":""}|로그인|
/usr/ins|POST|username=홍길동&nickname=dkagh&hp=01000000000&email=test@test.com&nickname=닉네임&pwd=dkagh| - |{"code": 200, "message":""}|회원가입|
/usr/get|GET| - | - |{"code": 200, "message":"", "username":"홍길동","nickname":"dkagh","hp":"01000000000","email":"test@test.com","nickname":"닉네임"}|회원정보조회|
/usr/dup/:usertype/:typevalue|GET| {"Hp"}/{"01000000000"} | - |{"code": 200, "message":""}|중복체크|
/usr/upd|PATCH|username=홍길동&nickname=dkagh&hp=01000000000&email=test@test.com&nickname=닉네임&pwd=dkagh| - |{"code": 200, "message":""}|회원정보 수정|
/usr/upd/pwd|PATCH|pwd=dkagh| - |{"code": 200, "message":""}|회원정보 수정|

DB & TABLES

* TBL_ADDITIONAL_USER 아이디로 사용할수 있는 필수정보 저장
* TBL_USER 유저고유키 이름 저장
* TBL_SCTRAN 문자인증관련 테이블
* TBL_TOKEN 토큰정보 저장

* 컬럼설명은 DDL에 저장

특히신경쓴부분
* 암호화
* 닉네임,이메일,휴대폰번호로 로그인 가능
* 휴대폰번호 가입시 기존에 동일한 휴대폰번호로 인증된 데이타가 있을 경우 처리방법 고민
* 중복된 데이타 처리방법 고민
