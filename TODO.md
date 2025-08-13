## Done
✅ - Feed endpoint <br />
✅ - Comments on posts
✅ - make possible to request posts by scrolling, aka from some index onwards, aka skip & take

## In Progress
🔶 - 

## Features
🕙 - On createPost, return entire post with full user and stuff, so that we dont have to fetch all the posts again, just to display the last one with author name.

## Architecture
🕙 - get rid of repeated logic for including the comments in getFeed and GetPostsByUser<br />
🕙 - implement 3-layered architecture<br />
🕙 - Migrate to Postgres<br />
🕙 - implement db migrations

## Security
🕙 - limit size of every single inputtable field in order to prevent flooding attacks<br />
🕙 - rate-limit per user to mitigate brute-force attacks<br />
🕙 - two-factor auth for user security and to obstruct fake profile creation<br />
🕙 - CAPTCHA to obstruct fake profile creation
