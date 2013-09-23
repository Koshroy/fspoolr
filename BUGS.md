Bugs
====

[checkmark] a build command changes the current directory and the old directory is not reset
* fsnotify spawn multiple events for single actual filesystem event, so do not run build command more than once
* does not detect all events