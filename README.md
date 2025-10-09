just project for fun which is basically identical functionally to what <a href="https://github.com/ThePrimeagen" target="_blank">ThePrimeagen</a> did on stream. (but it's still raw)

<strong>Ease vote</strong> allows you to write comments, modify, delete and search through database.

<strong>register/login</strong> is handled with <strong>jwt</strong> instead of twitch oauth like on stream.

just like the original one, comments are <strong>exclusive</strong>, meaning as an user you can post only 1 comment. 

app has <strong>like/dislike</strong> counts which sorts comments and shows the most liked ones first.

app is completely rendered on server, therefore it uses Templ library along with htmx to render front. <br> another difference is that <strong>Easy vote</strong> is powered by Golang on Echo framework (original is on php/laravel)

project is mostly for portfolio and is not meant to be used in production (even tho with slight modifications it can work there too without any problems)












