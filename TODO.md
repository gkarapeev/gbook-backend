## Done
✅ - Feed endpoint<br />
✅ - Comments on posts<br />
✅ - make possible to request posts by scrolling, aka from some index onwards, aka skip & take<br />
✅ - Make a get-user request that reutrns the full name and stuff by id<br />
✅ - get rid of repeated logic for including the comments in getFeed and GetPostsByUser<br />
✅ - Migrate to Postgres<br />
✅ - Also enable TLS on test
✅ - implement db migrations

🔶 - Add migrations for creating the basic tables - users, etc.

## Architecture
🕙 - implement 3-layered architecture<br />

## Security
🕙 - limit size of every single inputtable field in order to prevent flooding attacks<br />
🕙 - rate-limit per user to mitigate brute-force attacks<br />
🕙 - two-factor auth for user security and to obstruct fake profile creation<br />
🕙 - CAPTCHA to obstruct fake profile creation
