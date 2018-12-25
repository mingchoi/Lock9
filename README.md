# Lock 9

A telegram bot provides vote feature across multiple groups.

[![Build Status](http://drone.fallen.world/api/badges/mingchoi/Lock9/status.svg)](http://drone.fallen.world/mingchoi/Lock9)

## How to use

Add `@lock9_bot` to your group, then type one of the following command:

### Vote

```
# Start a vote quickly
/vote Topic OptionA OptionB

# Start a vote with options
/voteadv {single|multiple} Topic OptionA OptionB

# Forward a vote
/forwardvote {VoteID}

```

### Accounting

```
# Transfer money to someone by:
/atm Title @payer 1500yen @payee

# Lend money to peoples by:
/lend Title @payer @lenderA 1500yen @lenderB 1800yen...

# Split a bill to all people by:
/split Title @payer 3000 @lenderA @lenderB...

```

### Change Log

v0.2 - Accounting feature added
v0.1 - Vote feature added
