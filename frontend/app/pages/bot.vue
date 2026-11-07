<script setup lang="ts">
const ApiSniper = async () => {
    const res = await fetch("https://member.lazada.co.th/user/api/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "include", // สำคัญ: คำสั่งนี้จะบอกให้เบราว์เซอร์ดึง Cookies ทั้งหมด (รวม HttpOnly) ส่งไปด้วยอัตโนมัติ
        body: JSON.stringify({
                            "module": {
                                "sessionId": "13173b0845fc1bbadb8486e45f4b2a00",
                                "sessionExpiredTime": 1794223722905,
                                "refreshToken": null,
                                "tokenExpiredTime": null,
                                "cookieList": [
                                    {
                                        "name": "lzd_uid",
                                        "value": "10883271",
                                        "comment": null,
                                        "domain": ".lazada.co.th",
                                        "maxAge": -1,
                                        "path": "/",
                                        "secure": false,
                                        "version": 0,
                                        "httpOnly": false,
                                        "sameSite": "None"
                                    },
                                    {
                                        "name": "lzd_sid",
                                        "value": "13173b0845fc1bbadb8486e45f4b2a00",
                                        "comment": null,
                                        "domain": ".lazada.co.th",
                                        "maxAge": -1,
                                        "path": "/",
                                        "secure": false,
                                        "version": 0,
                                        "httpOnly": true,
                                        "sameSite": "None"
                                    },
                                    {
                                        "name": "lzd_b_csg",
                                        "value": "c655e401",
                                        "comment": null,
                                        "domain": ".lazada.co.th",
                                        "maxAge": -1,
                                        "path": "/",
                                        "secure": false,
                                        "version": 0,
                                        "httpOnly": true,
                                        "sameSite": "None"
                                    },
                                    {
                                        "name": "sgcookie",
                                        "value": "E100gU7DX7zIC95UtUnPnxD35obklaP85X8e2ACL9cSS0esZX2GKKGqsvAgMS9umAJRfPkFnxIelphOaniRSWTisjs3x2%2B%2F98uvJPDgT4TzAoN8%3D",
                                        "comment": null,
                                        "domain": ".lazada.co.th",
                                        "maxAge": 2592000,
                                        "path": "/",
                                        "secure": false,
                                        "version": 0,
                                        "httpOnly": true,
                                        "sameSite": "None"
                                    }
                                ],
                                "cookies": [
                                    "lzd_uid=10883271; Domain=.lazada.co.th; Path=/",
                                    "lzd_sid=13173b0845fc1bbadb8486e45f4b2a00; Domain=.lazada.co.th; Path=/; HttpOnly",
                                    "lzd_b_csg=c655e401; Domain=.lazada.co.th; Path=/; HttpOnly",
                                    "sgcookie=E100gU7DX7zIC95UtUnPnxD35obklaP85X8e2ACL9cSS0esZX2GKKGqsvAgMS9umAJRfPkFnxIelphOaniRSWTisjs3x2%2B%2F98uvJPDgT4TzAoN8%3D; Domain=.lazada.co.th; Expires=Wed, 02-Dec-2026 11:28:42 GMT; Path=/; HttpOnly"
                                ],
                                "headers": {
                                    "s_v": "4.0.2.6",
                                    "s_tag": "39582422106112|0^1|0^lazadaBuyer:10883271:^lazadaBuyer#8589967360|0",
                                    "s_write_units": "[SoutheastAsia:CENTER]",
                                    "s_ucode": "SoutheastAsia:CENTER",
                                    "s_tid": "2141118417936189227775252e6369",
                                    "s_read_unit": "[SoutheastAsia:CENTER]",
                                    "s_group": "lazada-buyer-session_LAZADA_TH",
                                    "s_status": "STATUS_NORMAL"
                                },
                                "loginType": "EMAIL",
                                "isAutoRegister": null,
                                "user": {
                                    "id": 10883271,
                                    "name": "Akkarapon",
                                    "birthday": "2004-03-28",
                                    "email": "qwasdzx225@gmail.com",
                                    "emailConfirmed": false,
                                    "phone": "0647507253",
                                    "phonePrefixCode": "66",
                                    "gender": "male",
                                    "avatar": "https://lh3.googleusercontent.com/a-/AOh14Gi-JyI42evLu2TPCUuS7EVT5Aae2I1NPBPT7kK6-Q=s96-c",
                                    "passwordVersion": null,
                                    "alipayUserId": "2165210011869114",
                                    "language": "1",
                                    "taxId": null,
                                    "branchId": null,
                                    "contactEmail": "qwasdzx225@gmail.com",
                                    "contactPhone": null,
                                    "isVerified": false,
                                    "b2b": false,
                                    "isTest": false,
                                    "status": "NORMAL",
                                    "type": "BUYER",
                                    "isLiveUp": false,
                                    "liveUpStatus": "NO_MEMBER_FOR_NEW",
                                    "enableEwallet": true,
                                    "hasAddress": null,
                                    "memberLevel": "TRUSTED",
                                    "enableNewsletter": false,
                                    "enableSmsNewsletter": false,
                                    "userTag1": null,
                                    "userTag2": null,
                                    "userTag3": null,
                                    "userTag4": null,
                                    "userTag5": null,
                                    "ext": "{\"emailConfirmed\":false,\"language\":\"1\",\"isAdmin\":true,\"havanaId\":614183012544}",
                                    "gmtCreate": null,
                                    "gmtModified": null,
                                    "siteId": "TH",
                                    "hasPassword": true,
                                    "hasPinCode": null,
                                    "isShowPinCode": null,
                                    "isShowBiometric": null,
                                    "site": 32,
                                    "originalPhone": "647507253",
                                    "adminUserId": null,
                                    "enableQuickLogin": true
                                },
                                "sgcookie": "E100gU7DX7zIC95UtUnPnxD35obklaP85X8e2ACL9cSS0esZX2GKKGqsvAgMS9umAJRfPkFnxIelphOaniRSWTisjs3x2%2B%2F98uvJPDgT4TzAoN8%3D",
                                "loginNextStepAction": null,
                                "loginNextStepActionData": null,
                                "isEmailOptional": null,
                                "isExclusiveFan": null,
                                "isShopOwner": null,
                                "loginSuccessTrustToken": null,
                                "oauthType": null,
                                "loginBehavior": "EMAIL_LOGIN",
                                "secondVerifyType": null,
                                "allLoginMethods": null,
                                "hasBiometric": false,
                                "supportThirdType": [
                                    "GOOGLE"
                                ],
                                "lazadaId": 999963143682,
                                "deviceTokenDTO": null,
                                "quickToken": null,
                                "quickTokenExpireTime": null,
                                "userId": 10883271,
                                "name": null,
                                "avatar": null,
                                "needCompletePassword": false,
                                "completePasswordToken": null,
                                "completeEmailToken": null,
                                "toRegisterToken": null,
                                "userVerification": null,
                                "hasPhone": true,
                                "hasEmail": true,
                                "oauthUserId": null,
                                "oauthUserEmail": null,
                                "oauthUserName": null,
                                "isRefreshSession": null,
                                "redirectUrl": null,
                                "errorCode": null,
                                "loginSuccess": true
                            },
                            "success": true,
                            "errorCode": null,
                            "retry": false,
                            "repeated": false,
                            "notSuccess": false
                        })
    })
    const rawText = await res.text();
    console.log("Raw Response:", rawText);
    try {
    // ลองแปลงเป็น JSON ถ้าทำได้
    const data = JSON.parse(rawText);
    console.log("JSON Data:", data);
    } catch (e) {
    console.warn("เซิร์ฟเวอร์ไม่ได้ตอบกลับมาเป็น JSON (อาจส่งหน้า HTML / Captcha มาแทน)");
    }
};
</script>

<template>
    <div>
    <button type="button" @click="ApiSniper">
        กด
    </button>
    </div>
</template>
