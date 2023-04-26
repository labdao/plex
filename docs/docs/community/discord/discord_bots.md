---
title: Research Assistant Bot
description: Research Assistant is a general purpose discord bot made custom for the LabDAO discord server.
sidebar_position: 2
---

## Introduction 

Research Assistant is a general purpose discord bot made custom for the LabDAO discord server. The bot is in active development and new functions are consistently being added when new needs arise in the discord. Please comment on Github if you have an idea you believe improve this discord.

## Command Prefix and Docs Syntax

### Docs Syntax
* ```<command>``` indicates command calls the function. 
* ```{argument}``` indicates the argument is required for the function. 
* ```[argument]``` indicates the argument is optional.
* ```[argument]*(x -> y)``` indicates the command can take between x and y optional arguments. 
* ```<command>\*\*``` indicates that the command requires specific permissions 


### Calling Commands
* Command Prefix\: $
    * ```$<command> [argument]``` calls a command
* Mentioning Bot will also call commands
    *  ```@Research Assistant <command> [argument]```
* The bot is agnostic to capitalization on command calls but is required to match in arguments.

## Current Live Commands

### Ping 
Ping checks the bot is responding and provides the current latency. 

* ```\$\<ping\>:```


![](https://i.imgur.com/edjIAKD.png)

### Link

Link return server relevant links in the form of a discord embed. 

* ```\$\<link> [website name]```
* website name is unformatted text that must match one of the links provided in the screen shot to get a result. 
* The Bot has a programmed error message if there is no match.

If no argument is provided, the bot returns a bank of all URLs hyperlinked in an embed. 

![](https://i.imgur.com/nEoJiWs.png)



If an argument is provided, the bot returns the link of that type: 

![](https://i.imgur.com/R21F9vI.png)

At the moment, there are two single word commands that the bot will recognize. 

```\$github``` and ```\$twitter``` will return the links to the github and the twitter account respectively.


### Poll 
 
Poll creates an interactive embed based poll where users vote via emoji reacts. 

```
\$\<Poll> {"Question"} ["Option"]*(0->6)
```

All arguments in this function must be contained in quotation marks for proper bot parsing, and are otherwise text strings.

No arguments provided after a question will default to Yes and No as two options to vote, however providing a question is mandatory. 

Providing multiple arguments will insert those arguments into the poll and have an emoji vote for each one.

![](https://i.imgur.com/Nb3nL0N.png)

### Help

A help function that returns a list of current active functions. 

```\<help\> [command]```

Returns a message with a list of all the function names and function description subject to character limits.

If optional command argument is provided, the help function returns information on the provided function. 

### Admin and Moderation\*\* 

For these commands to recognize users, the name must match exactly. It is recomended you copy and paste their discord user name. 
The four digit descriminator after the # is not needed. Nicknames may not be recognized by the bot.
Role title must also match exactly, including emojis if they are included. Copy-Paste Recomended.

* Role (requires manage role permissions)
    * ```\<Role\> {User} {Role}```
        * Gives user role if they do not have it 
        * Removes role if the user has it. 

* Ban
    * ```\<Ban\> {User}```
    * Bans user
* Hackban
    * ```\<Hackban\> {User}```
    * Bans user without them needing to be in the server
* Kick
    * ```\<Kick\> {User}```
    * Removes a user from your server. 
* Purge```
    * ```\<Purge\> {Integer}```
    * Removes previous messages until {Integer} messages have been removed

## Backend Functions 

* Reaction Roles

![](https://i.imgur.com/CADc4dq.png)

Bot can be configured to automatically assign roles to users upon reacting to a message. 
Can be configured to respond to any number of messages with various sets of emojis and roles. 

At the time these are programmed manually and upon request, but creating a discord text-based UI for creating new reaction roles can be done.

* New User Welcome Message and DM (Undeployed)
    * Code exists for the bot to automatically dm and send a welcome message to new users when they join the server. 
    * Messages for this is entirely customizable. 


### "Fun"ction and Methods 

* GM - Says a custom or random message once per day in the #gm channel. 
* STATUS - Changes self status automatically after a set period of time. 
* ```\<8ball\> {question}```, the bot will answer your question as if it was a magic 8ball.

