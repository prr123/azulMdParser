# azulMdParser

## azmd Format
azmd is an extended markdown format. 
The format has three parts:
 1. A heading section encapsulated by ---/n. 
The heading section contains information in the yaml format.
Sample Format:  
`
title: "Title Test: The Son Rises"
author: Peter Riemenschneider
short: prr
date: 20/8/2025
keys:
- "#Valencia"
- "#spacesCabo"
- "#homes"
`
2. A summary section.
The summary section is a description of the article in the markdown format.
The summary has to start with a line "#Summary" and ends with the next line starting with a hash tag (#).  

3. The main section
The main section is an article in the markdown format.

## the Parser
The parser reads a file from azmd directory and assumes an "azmd" extension.  
The parser generates a sub directory in the directory named "articles".  

The parser will write three files in the sub directory:  
- a header file with a yaml extension
- a md file with the name summary.md
- a md file with the name main.md



## Sample File
```
---
title: "Title Test: The Son Rises"
author: Peter Riemenschneider
short: prr
date: 20/8/2025
keys:
- "#Valencia"
- "#spacesCabo"
- "#homes"
---
# Summary
This is a summary of the article. It can be up to 500 chars long

# Article
## Introduction
hello
## The Market
the market is changing.
## Conclusion
The north is still interesting.
```

## directories

articles azmd

|---------|
|  azmd   | 
|---------|

reads and azmd file from the azmd directory

the azmdMdParser generates three files in the articles directory

|------------------------------|
| head.yaml summary.md main.md |
|------------------------------|

