function checkPhone(value) {
    const numberreg = /^((0\d{2,3}-\d{7,8})|(1[3456789]\d{9}))$/ // 手机号
    return numberreg.test(value)
}

function trim(text) {
    return text.replace(/(^[\s]+)|([\s]+$)/g, '');
}

function processImage(srcfile, quality, callback) {
    console.log('processImage quality', quality)
    const canvas = document.createElement('canvas');
    const context = canvas.getContext('2d');
    const reader = new FileReader();
    const img = new Image();
    img.onload = function () {
        const targetWidth = Math.round(this.width * 0.5);
        const targetHeight = Math.round(this.height * 0.5);
        canvas.width = targetWidth;
        canvas.height = targetHeight;
        context.clearRect(0, 0, targetWidth, targetHeight);
        context.drawImage(img, 0, 0, targetWidth, targetHeight);
        canvas.toBlob(blob => {
            blob.name = srcfile.name;
            callback(blob);
        }, srcfile.type, quality);
    };
    reader.onload = e => {
        img.src = e.target.result;
    };
    reader.readAsDataURL(srcfile);
}

function getBody(xhr) {
    const text = xhr.responseText || xhr.response;
    if (!text) {
        return text;
    }
    try {
        return JSON.parse(text);
    } catch (e) {
        return text;
    }
}

function getErrtext(xhr) {
    if (xhr.response) {
        return xhr.status + ":" + (xhr.response.error || xhr.response)
    } else if (xhr.responseText) {
        return xhr.status + ":" + xhr.responseText
    } else {
        return "fail to post " + xhr.status;
    }
}

function uploadpic(file, errgo, sideinfo) {
    var size = file.size / 1024 / 1024
    console.log("uploadpic size",size)
    if (size > 2) {
        errgo.html('上传文件太大，超出2M限制') 
        return false
    }
    var upurl = rootPath + "/uploadidnum"
    console.log("url", upurl)
   
    var xhr = new XMLHttpRequest();
    xhr.open('POST', upurl, true);
    // xhr.setRequestHeader("Content-type", "multipart/form-data;");//不用设
    const formData = new FormData();
    formData.append("upfile", file);
    formData.append("side", sideinfo);
    loadinggo.show()
    xhr.onerror = function error(err) {
        loadinggo.hide()
        console.log("errr", err)
        errgo.html("网络错误")
    };

    xhr.onload = function onload() {
        loadinggo.hide()
        if (xhr.status < 200 || xhr.status >= 300) {
            console.log("err", getErrtext(xhr))
        }
        var getres=getBody(xhr)
        if(getres.code==1){
            // alert("上传成功")
            errgo.html("上传成功")
            var getdata=getres.data
            if(sideinfo=="front"){
                idpicitem1.attr("src", getres.data.url)
                idcardnum.val(getres.data.result["公民身份号码"].Words)
                idcardname.val(getres.data.result["姓名"].Words)
            }else{
                idpicitem2.attr("src", getres.data.url)
            }

        }else{
            errgo.html(getres.message)
        }
        console.log("get ok", getBody(xhr))
    };

    if (xhr.upload) {
        xhr.upload.onprogress = function progress(e) {
            if (e.total > 0) {
                loadinggo.show()
                progressStatus.html('图片上传中') ;
            }
            console.log("progress", e)
        };
    }


    xhr.send(formData);

}

function updloadfile(file, errgo, sideinfo) {
    
    console.log('begin uploadprocess', file)
    if(!file){
        errgo.html(' 没有选择文件') 
        return false
    }
    if (file.type.indexOf('image') === -1) {
        errgo.html('只接受png和jpeg格式') 
        return false
    }

    var size = file.size / 1024 / 1024
    console.log('origin size', size)
    if (size > 1) {
        console.log('begin compress')
        processImage(file, 1, function (bob) {
            uploadpic(bob, errgo, sideinfo)
        })
    } else {
        console.log('upload')
        uploadpic(file, errgo, sideinfo)
    }
}
var idpicerr2 = $("#back-side-warning")
var idpicerr1 = $("#front-side-warning")
var phoneerr = $("#mobile-error")
var phoneinput = $('input[name=mobile]')
var idpicitem1 = $("#idpic1")
var idpicitem2 = $("#idpic2")
var idcardnum = $('#idcardnum')
var idcardname = $('#idcardname')
var loadinggo = $("#loading")
var progressNumber = $("#progressNumber")
var progressStatus = $("#progressStatus")
var submiterrGo=$("#submit-error")
var submitokGo=$("#submit-success")
var curPageUrl = window.document.location.href;
var picfrontsrc="/static/images/desktop-front-side.png"
var picbacksrc="/static/images/desktop-back-side.png"
var rootPath = curPageUrl.split("//")[0] + '//' + curPageUrl.split("//")[1].split("/")[0]
console.log("rootpath", rootPath)

function emptyall(){
    idcardname.val('');
    idcardnum.val('');
    phoneinput.val('');
    phoneerr.empty()
    idpicerr2.empty()
    idpicerr1.empty()
    loadinggo.hide()
    submiterrGo.empty()
    submitokGo.empty()
    idpicitem1.attr("src", picfrontsrc)
    idpicitem2.attr("src", picbacksrc)
}
//验证
$(document).ready(function (e) {

    console.log("test ready",phoneerr)
    emptyall()

    $("#idpicinput1").change(function (event) {
        console.log("el", event.target)
        idpicerr1.empty()
        var fileinfo = event.target.files[0];
        console.log("fileinfo", fileinfo)
        updloadfile(fileinfo, idpicerr1, "front")
    })

    $("#idpicinput2").change(function (el) {
        console.log("el", event.target)
        idpicerr2.empty()
        var fileinfo = event.target.files[0];
        console.log("fileinfo", fileinfo)
        updloadfile(fileinfo, idpicerr2, "back")
    })

    var checkphoneexe=function(){
        phoneerr.empty()
        var mobiletext =phoneinput.val();
        if (checkPhone(trim(mobiletext)) == false) {
            phoneerr.html("不是正确的手机号") 
            return false
        }
        return true
    }
    phoneinput.keyup(function () {
        checkphoneexe()
    })

    $("#submit").click(function(el){
        console.log("click")
        submiterrGo.empty()
        submitokGo.empty()
        var mobiletext =phoneinput.val();
        if (checkPhone(trim(mobiletext)) == false) {
            submiterrGo.html("不是正确的手机号") 
            return false
        }
        var idpic1=idpicitem1.attr("src")
        if (idpic1===""||idpic1===picfrontsrc) {
            submiterrGo.html("身份证正面没有上传") 
            //idpic1="";
            return false
        }
        var idpic2=idpicitem2.attr("src")
        if (idpic2===""||idpic2===picbacksrc) {
            submiterrGo.html("身份证反面没有上传") 
            //idpic2="";
            return false
        }
        var cardnum=trim(idcardnum.val())
        if (cardnum==="") {
            submiterrGo.html("身份证号空") 
            return false
        }
        var cardname=trim(idcardname.val())
        if (cardname==="") {
            submiterrGo.html("收件人姓名空") 
            return false
        }

        var updatainfourl = rootPath + "/Logistics/ClientChangeInfo"
        console.log("updatainfourl", updatainfourl)
        loadinggo.show()
        $.ajax({
            url: updatainfourl, //百度接口api 鹰眼
            type: 'POST', //GET
            contentType:"application/json;charset=utf-8",
            async: true, //或false,是否异步
            data: JSON.stringify({idnumpic1: idpic1,
                idnumpic2: idpic2,
                idnum: cardnum,
                client_name: cardname,
                client_phone: mobiletext}),
            timeout: 90000, //超时时间
            dataType: 'json', //返回的数据格式：json/xml/html/script/jsonp/text
            beforeSend: function (xhr) {
                console.log(xhr)
                console.log('发送前')
            },
            success: function (res, textStatus, jqXHR) {
                console.log('res',res)
                loadinggo.hide();
                if(res.code==1){
                    submiterrGo.empty()
                    submitokGo.html("上传成功")
                    alert("上传成功！");
                }else{
                    submiterrGo.html(res.message)
                }
            },
            error:function(xhr,errstatus,err){
                loadinggo.hide();
                console.log('err',errstatus,err)
                submiterrGo.html("网络错误")
            },
            complete:function(xhr,status){
                loadinggo.hide();
                console.log('complete',xhr,status)
            }
        })

    })

});