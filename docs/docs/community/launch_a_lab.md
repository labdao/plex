---
title: Setting up a project
description: To collaborate effectively online, we use community standards to organise scientific results, code and writings.
sidebar_position: 4
---

For your scientific work to be reproducible and collaborative, we recommend going through a set of steps when launching: 
1. set up group chats for your project
2. integrating with the LabDAO discord to increase visibility
3. set up a GitHub repository for your project
4. Get started with a cookiecutter directory structure for your code
5. request a Google and Notion Workspace for the team

## Setting up group chat channels
We generally recommend you set up a (1) internal and (2) public facing communication channel for your project. The internal channel is for committed contributors to the project, while the public channel is a place for you to share updates with a wider group of interested community members and enthusiasts.

Telegram is a free messaging service that enables groups to collaborate very easily. While it is not the most secure messaging service (we recommend [Signal](https://www.signal.org/) for this), it supports a lot of bots which can be useful to easily scale up a conversation and control access. For example, maintainers at LabDAO can token-gate access to telegram channels easily - with the right token, you can grant access to all relevant ressources of a [laboratory](https://guild.xyz/labdao) for easy onboarding of new contributors. 

## Integrating with the LabDAO discord to increase visibility
Next to an independent set of communication channels, maintainers within LabDAO can create an internal and public communication channel within the LabDAO discord, too. This increases the visibility for your project and unlocks the micro-expertise that might be available within the community. 

## Setting up a GitHub repository
If you want to maximise the visibility within the community, reach out to a steward and set it up in [LabDAO-Projects](https://github.com/labdao-projects).



## installing cookiecutter datascience and creating a directory

Cookiecutter datascience is a community standard for scientific project directories. By following an agreed upon pattern about what files to put where, we can accelerate the rate of collaboration between scientists.

````
# Python 2.7 or 3.5 required
pip3 install cookiecutter
cookiecutter https://github.com/drivendata/cookiecutter-data-science
# it does not hurt to reinstall the tool
````

## setting up a project repository
To set up a standardised directory structure, a set of inputs are needed during the setup:
1. *project_name* - Project names usually follow the pattern *labdao_yourname*. 
2. *description* - a summary of the research project you are conducting.
3. *license* - LabDAO :green_heart: open source. Sometimes, however, the research you are doing is processing critical health information (not recommended at this point!) or leading to valuable new insights that need to be protected in order to make it to market. In this case you should make sure that all data and code in your GitHub repository and S3 buckets are protected. If you are not sure what license is the best, choose "no license" for the start. 
4. *bucket* - while we are still building infrastructure for decentralised storage, we recommend using an AWS S3 bucket for remote file storage. Community members can help you set up an S3 bucket. 

## tracking your project directory with git
Once the project is set up, initialize a git repository by calling ```git init```. The directory already comes with a set of useful defaults on what files and directories not to include in your repository.

## sharing files 
Files located in ```/data``` can be shared with everyone in the team by uploading to the group's bucket: 

```
make sync_data_to_s3
```
on the receiving end data can be synced from the group's bucket using: 

```
make sync_data_from_s3
```



