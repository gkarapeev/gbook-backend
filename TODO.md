## Done
✅ - Feed endpoint <br />
✅ - Comments on posts

## In Progress
🔶 - 

## Architecture
🕙 - get rid of repeated logic for including the comments in getFeed and GetPostsByUser<br />
🕙 - implement 3-layered architecture<br />
🕙 - Migrate to Postgres<br />
🕙 - implement db migrations

## Performance
🕙 - make possible to request posts by scrolling, aka from some index onwards, aka skip & take

## Security
🕙 - limit size of every single inputtable field in order to prevent flooding attacks<br />
🕙 - rate-limit per user to mitigate brute-force attacks<br />
🕙 - two-factor auth for user security and to obstruct fake profile creation<br />
🕙 - CAPTCHA to obstruct fake profile creation
