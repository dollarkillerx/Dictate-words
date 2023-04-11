import Koa from 'koa'
import KoaRouter from 'koa-router'
import bodyParser from 'koa-bodyparser'
import * as googleTTS from 'google-tts-api';
import speak from 'google-translate-api-x';
import { writeFileSync } from 'fs';

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

    speak('gata', {to: 'es'})


    ctx.response.type = 'application/json';
    ctx.body = {
        "url": url
    };
})

router.post('/google_tts2', async (ctx) => {
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

    if (lang === 'jp' ) {
        lang = "ja"
    }


    // const res = await  speak(text, {to: lang})

    const res = await speak('gata', {to: 'es'}); // => Base64 encoded mp3
    // writeFileSync('cat.mp3', res, {encoding:'base64'}); // Saves the mp3 to file

    // writeFileSync('cat.mp3', res, {encoding:'base64'});

    ctx.response.type = 'application/json';
    ctx.body = {
        "base64": res
    };
})


async function main() {
    app.use(router.routes());
    app.listen(3030);
}

main().then(r => {
    console.log('server is running in', 3030)
})