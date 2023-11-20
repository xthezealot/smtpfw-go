# SMTP Forwarder

Listens and forwards SMTP messages to another sever.

It's mainly used as a bridge between old mail clients and mail servers that don't accept old TLS versions anymore (like AWS SES).

Edit the settings in `smtpfw.json`.

**WARNING:** This cannot be complied for Windows XP or older because dependencies don't support Go v1.10 (the last version supporting Windows XP) anymore. See the [Python 3.4 variant](https://github.com/xthezealot/smtpfw) instead.
