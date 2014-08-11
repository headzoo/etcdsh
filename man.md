ETCDSH 0.2 "AUGUST 2014" "etcd man page"
========================================

NAME
----
etcd

SYNOPSIS
--------
etcd [`help`]

OPTIONS
-------
`help` Display the help page and exit.

COMMANDS
--------
The following commands are available from within the shell.

`ls` - Lists the contents of the working directory. The output contains 5 columns: CreatedIndex, ModifiedIndex, TTL, Type, Name. For example:

total 4  
 0  0   0 k .  
 0  0   0 k ..  
12 12 300 k domains  
34 34   0 o version

Using the "domains" row, the first column is the created index (12). The second is the modified index (12). The third column is the TTL (300). The fourth column is either "k" or "o" (k). A "k" means the value in the fifth column is a key, and "o" means it's an object. The fifth column is the name of the key or object (domains).

Examples:
 ls  
 ls /  
 ls domains  
 ls ..

`cd` - Change the working directory.

Examples:
 cd /  
 cd /domains  
 cd ..  
 cd ../..  
 cd /domains/apps

`get` - Displays the value of an object.

Examples:
 get /domains/apps  
 get /versions/1.0

`set` - Sets the value of an object.

Examples:
 set /version/app 1.0  
 set /domains/apps/mobile mobile.xdt.io


AUTHOR
------
Sean Hickey (sean@headzoo.io)
