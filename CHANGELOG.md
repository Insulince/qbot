# CHANGELOG

### v0.8.2 - 16 May 2025
- `!history list` now supports offset
- `!history list` now returns detailed stats about each tourney, such as the winner, winner's waves, average wave count, and number of entrants

### v0.8.1 - 16 May 2025
- Add `list` mode to the `!history` command

### v0.8.0 - 15 May 2025
- Display names are used instead of usernames in the loser meme.
- Add some dev utilities

### v0.7.10 - 11 May 2025
- @everyone announcements now come with an image.

### v0.7.9 - 11 May 2025
- Freaking !owned.

### v0.7.8 - 11 May 2025
- Updated the last place meme to be dynamic, showing the username of the last place player regardless of who it is
- Added image manipulation to overlay the last place player's username directly onto the meme image
- Improved text rendering with white text and black outline/shadow for better visibility
- Right-aligned text in the lower right corner of the image for better aesthetics
- Replaced "wompwomp.png" with "celebrate.png" for the last place celebration image

### v0.7.7 - 21 March 2025
- !q is now a valid way to enter the queue
- if you put your command in all caps, Q will get mad at you

### v0.7.6 - 19 March 2025
- enter -> full timeout is now **30 minutes** (up from 15)

### v0.7.5 - 15 March 2025
- but not too funny

### v0.7.4 - 12 March 2025
- memes r funny

### v0.7.3 - 8 March 2025
- Setup role mention for error reporting to get my attention quicker.

### v0.7.2 - 8 March 2025
- Correct bug introduced by new command and argument processing logic that allows commands which requires arguments to function properly again (i.e. `!submitwave`)

### v0.7.1 - 8 March 2025
- Add `!deverror` for developer testing of error handling by the bot.

### v0.7.0 - 8 March 2025
- Overhaul error handling and a lot of architectural changes in general.
- Errors should now be posted to a private discord channel for easier and more timely review.

### v0.6.3 - 7 March 2025
- SQL query for inserting tourney entries had an error when checking CONFLICTS. This has been corrected!

### v0.6.2 - 7 March 2025
- `short_name` was being referred to as `shortName` in SQL queries, which is incorrect. This was causing a new tourney to not be created on tourney open.

### v0.6.1 - 7 March 2025
- `!submitwave` can also be called as `!submitwaves`

### v0.6.0 - 7 March 2025
- Create `!history` command for viewing historical tournament results
- Setup automatic tournament creation in sql when a new tourney opens via the scheduler
- `!leaderboard` will now pull from the latest tourney, to view previous ones use !history
- `!submitwave` also operates on the latest tourney.
- `!view` of the queue should no longer ping anyone from its tags.

### v0.5.1 - 5 March 2025
- Fix bug in scheduler so that announcements hit at the right time.

### v0.5.0 - 3 March 2025
- Add `!submitwave` command for submitting wave count. Right now its global and static, but eventually this will be per-tourney.
- Add `!leaderboard` command for displaying leaderboard of player submitted scores
- At tourney end, `!leaderboard` is automatically ran by Q to announce results.
- These will have much more work done to them, this is just a basic starting point to actually leverage the database.

### v0.4.0 - 3 March 2025
- A lot of tweaking of the Dockerfile to enable LiteFS for SQLite support in Fly.io.
- Add scheduler.go for sending out scheduled pings regarding tourney events.
- Add temporary `!insert`, `!fetch`, and `!deleteall` commands for testing the database. These will be deleted soon.
- Add disclaimer section to `!help`

### v0.3.0 - 28 February 2025
- Setup deployments to fly.io
- Setup autodeploy on push
- Add `!version` command

### v0.2.1 - 28 February 2025
- Setup justfile to allow for both mac and windows setups seamlessly
- Add git pre-commit hooks
- Create version.go for tracking versions

### v0.2.0 - 28 February 2025
- queue -> enter timeout is still 5 minutes, really this should be fine, if you aren't able to join a bracket 5 mins after becoming the active user, then you forfeit your turn (or can extend!)
- enter -> full timeout is now **15 minutes** (up from 10)
- timeout warnings now hit at the 2 minute mark instead of the 1 minute mark
- `!moretime` now works for both when you are waiting to enter a bracket and when you are waiting to mark your bracket full (previously only worked on the latter)
- `!enter` will now tell you if you aren't even in the queue yet and recommend `!queue` instead of just saying "its not your turn"

### v0.1.0 - 24 February 2025
- Initial bot created.
