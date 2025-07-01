# HyveKeeper

The goal of this project is to create a modern web panel for [Bhyve](https://bhyve.org/). 

# Setup
## Permissions
It is required that `rctl -u` can be run by the user that is running the backend.
One way to allow this is to add the following to a file in `/usr/local/etc/sudoers.d/`
```sh
ALL ALL=(ALL) NOPASSWD: /usr/bin/rctl -u *
