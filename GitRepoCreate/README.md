# GOGIT

## Installation and Usage 

### Set ENV Variable 
This tool requires an environment variable called 'GOGIT' to be set. This variable should contain your personal access token from Github. 

### Options

This tool accepts a few different command lime options. Only the -n,name flag is mandatory to initialize the repository and create the Github repo. 

-n,name  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Name of the repository that will be created. 

-e,email  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Your Github email address

-h,help  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Displays the help menu

-k,sshkey  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Specify an SSH key for your git repo 

-u,username  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Your Github username

-v,version
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Check the version of Gogit

### Examples 

```
gogit -n RepoName -e example@example.com -u examplename
```


