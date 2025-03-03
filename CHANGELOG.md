# CHANGELOG

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
