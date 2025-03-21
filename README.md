[![CI/CD Pipeline](https://github.com/flintdemon/english_island_bot/actions/workflows/ci.yml/badge.svg)](https://github.com/flintdemon/english_island_bot/actions/workflows/ci.yml)

This bot was written for an English language school. The bot offers to take a language test, gives your level and offers to send contacts to the school to sign up for a trial lesson.

The list of questions for the test is stored in a simple yaml file. The telegram token and chat id where the user's contact will be sent is defined in environment variables. The bot cen be packaged in docker. 

Variables: 

TELETOKEN - Telegram token of your bot
ADMIN_CHAT_ID - id of user to send test results and contact
