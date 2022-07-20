import Koa from 'koa'
import KoaRouter from 'koa-router'
import bodyParser from 'koa-bodyparser'
import * as googleTTS from 'google-tts-api';

const app = new Koa();

app.use(bodyParser());

const router = new KoaRouter()

router.get('/', async (ctx) => {
    ctx.body = "Dictate words (Implemented through google tts)  https://github.com/dollarkillerx/Dictate-words"
})

router.post('/google_tts', async (ctx) => {
    let lang = ctx.request.body.lang;
    let text = ctx.request.body.text;
    if (lang === undefined || text === undefined) {
        ctx.status = 400
        ctx.body = "なに"
        return
    }

    lang = lang.replace(/^\s*|\s*$/g,"");
    text = text.replace(/^\s*|\s*$/g,"");
    if (lang==="" ||text==="") {
        ctx.status = 400
        ctx.body = "なに"
        return
    }

    const url = googleTTS.getAudioUrl(text, {
        lang: lang,
        slow: false,
        host: 'https://translate.google.com',
    });


    ctx.response.type = 'application/json';
    ctx.body = {
        "url": url
    };
})

async function main() {
    app.use(router.routes());
    app.listen(3030);
}

main().then(r => {
    console.log('server is running in', 3030)
})