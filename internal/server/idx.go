package server

import "github.com/gin-gonic/gin"

func (s *Server) showIdx(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")
	ctx.Writer.Write([]byte(idx))
}

var idx = `
<!doctype html>
<html lang="en">
<head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://cdn.jsdelivr.net/npm/vue@2.7.0/dist/vue.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/axios@0.27.2/dist/axios.min.js"></script>
    <!-- Bootstrap CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-F3w7mX95PdgyTmZZMECAngseQB83DfGTowi0iMjiWaeVhAn4FJkqJByhZMI3AhiU" crossorigin="anonymous">

    <title>Dictate Words  外語聽寫</title>

    <style type="text/css">
        a{
            text-decoration: none; /* 去除默认的下划线 */
            color: #000;    /* 去除默认的颜色和点击后变化的颜色 */
        }
    </style>
</head>
<body>


<div id="app">
    <div class="container">
        <h1>Dictate Words  外語聽寫</h1>
        <br>
   		<span>服務器配置弱, 如果發生錯誤 請等待3~5分鐘</span>
  		<a href="https://github.com/dollarkillerx/Dictate-words" class="btn btn-primary">給作者一個Start https://github.com/dollarkillerx/Dictate-words</a>
        <a href="https://github.com/dollarkillerx/Dictate-words/issues">反饋意見</a>
        <br>

        <button type="button" class="btn btn-primary" @click="genTts">生成聽寫音頻</button>
        <a type="button" class="btn btn-success" :href="downPath" v-show="downPathShow" target="_blank">點我下載音頻</a>

        <br>
        <br>
        <span>選擇語言: </span>
        <select class="form-select" aria-label="Default select example" v-model="langOption" @change="changeLangOption($event)">
            <option value="ja">日語</option>
            <option value="ko">韓語</option>
            <option value="de">德語</option>
            <option value="fr">法語</option>
            <option value="es">西班牙语</option>
            <option value="pt">葡萄牙語</option>
            <option value="ru">俄語</option>
            <option value="en">英語</option>
            <option value="zh-CN">中文</option>
        </select>

		<br>
        <span>播放順序: </span>
        <select class="form-select" aria-label="Default select example" v-model="playOrderOption" @change="changePlayOrderOption($event)">
            <option selected value="default">循序播放</option>
            <option value="random">亂序播放</option>
        </select>
        <br>
        <span>綫路: </span>
        <select class="form-select" aria-label="Default select example" v-model="lineOption" @change="changeLineOptionOption($event)">
            <option selected value="default">默認</option>
            <option value="spare">備用</option>
        </select>

        <br>

        <span>重複次數:  (最大3次)</span>
        <input type="text"  class="form-control"  v-model="repeatTimes" placeholder="重複次數 最大3次"  type="number">
        <br>
        <label class="form-label">單詞 Or 語句: (單行最長200字, 每次生成最大100行, 單詞或句子 一詞一行 一句一行)</label>
        <textarea type="text"  class="form-control"  v-model="inputSearch" style="min-height: 600px"></textarea>
        <br>
    </div>


</div>

<script>
    //在页面未加载完毕之前显示的loading Html自定义内容
    var _LoadingHtml = '<div id="loadingDiv" style="display: none; "><div id="over" style=" position: absolute;top: 0;left: 0; width: 100%;height: 100%; background-color: #f5f5f5;opacity:0.5;z-index: 1000;"></div><div id="layout" style="position: absolute;top: 40%; left: 40%;width: 20%; height: 20%;  z-index: 1001;text-align:center;">生成音頻中請等待...  </div></div>';
    //呈现loading效果
    document.write(_LoadingHtml);

    //移除loading效果
    function completeLoading() {
        document.getElementById("loadingDiv").style.display="none";
    }
    //展示loading效果
    function showLoading()
    {
        document.getElementById("loadingDiv").style.display="block";
    }

    var app = new Vue({
        el: '#app',
        data: {
            repeatTimes: 3,
            langOption: '',
            inputSearch: '',
       		downPath: "",
            downPathShow: false,
            playOrderOption: 'default',
            lineOption: 'default',
        },
        mounted() {

        },
        methods: {
            changeLangOption(event) {
                console.log(event.target.value); // 打印的结果就是，我们选中的option里面的value值
                this.langOption = event.target.value;
            },
 			changePlayOrderOption(event) {
                console.log(event.target.value); // 打印的结果就是，我们选中的option里面的value值
                this.playOrderOption = event.target.value;
            },
    		changeLineOptionOption(event){
                console.log(event.target.value);
                this.lineOption = event.target.value
            },
            genTts() {
                console.log(this.langOption)
                console.log(this.inputSearch)
                this.inputSearch = this.inputSearch.replace(/^\s*|\s*$/g,"")

                this.repeatTimes = parseInt(this.repeatTimes)
                if (this.langOption==="default"||this.inputSearch == "" || this.repeatTimes <= 0 || this.repeatTimes > 3) {
                    alert("請認證填寫請求參數")
                    return
                }

                showLoading()

     			let spare = false
                if (this.lineOption === "spare") {
                    spare = true
                }

                axios.post('/generate_tts',{
                    "lang": this.langOption,
                    "text": this.inputSearch,
                    "repeat_times": this.repeatTimes,
					"play_order": this.playOrderOption,
   					"spare": spare,
                })
                    .then( (response)=> {
                        completeLoading()
                        console.log(response.data.id)
                        // window.open("/download_tts/" + response.data.id, '_blank');
						this.downPath = "/download_tts/" + response.data.id
                        this.downPathShow = true
 						if (response.data.word !== "") {
                            this.inputSearch = response.data.word
                        }
                    })
                    .catch(function (error) {
                        completeLoading()
                        console.log(error);
                        alert("錯誤: " + error.response.data)
                    });
            }
        }
    })
</script>

<!-- Option 1: Bootstrap Bundle with Popper -->
<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.1/dist/js/bootstrap.bundle.min.js" integrity="sha384-/bQdsTh/da6pkI1MST/rWKFNjaCP5gBSY4sEBT38Q/9RBh9AH40zEOg7Hlq2THRZ" crossorigin="anonymous"></script>

<!-- Global site tag (gtag.js) - Google Analytics -->
<script async src="https://www.googletagmanager.com/gtag/js?id=UA-110586888-1"></script>
<script>
    window.dataLayer = window.dataLayer || [];
    function gtag(){dataLayer.push(arguments);}
    gtag('js', new Date());

    gtag('config', 'UA-110586888-1');
</script>


</body>
</html>
`
