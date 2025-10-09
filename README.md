just project for fun which is basically identical functionally to what <a href="https://github.com/ThePrimeagen" target="_blank">ThePrimeagen</a> did on stream. (but it's still raw)

<strong>Ease vote</strong> allows you to write comments, modify, delete and search through database.

<strong>register/login</strong> is handled with <strong>jwt</strong> instead of twitch oauth like on stream.

just like the original one, comments are <strong>exclusive</strong>, meaning as an user you can post only 1 comment. 

app has <strong>like/dislike</strong> counts which sorts comments and shows the most liked ones first.

app is completely rendered on server, therefore it uses Templ library along with htmx to render front. <br> another difference is that <strong>Easy vote</strong> is powered by Golang on Echo framework (original is on php/laravel)

project is mostly for portfolio and is not meant to be used in production (even tho with slight modifications it can work there too without any problems)


<img width="1920" height="953" alt="screencapture-localhost-8080-v1-comments-2025-10-09-18_00_00" src="https://github.com/user-attachments/assets/2f135940-c1ae-4863-b65a-7f55d5e010e4" />
<img width="452" height="958" alt="image" src="https://github.com/user-attachments/assets/5e3588a8-1df5-4cfe-b62d-a2df62879142" />







