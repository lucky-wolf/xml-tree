# xmltree
XML loading, saving, and manipulation at a generic level

Warn: sadly, golang has never updated their xml parser to allow for version 1.1, it only handles 1.0.  In my current experience, most files can still be read as 1.0 - 1.1 mostly adds a few otherwise restricted characters, which for English speakers, is unlikely to occur.  Ideally, Golang would update their code, but so far (2024-11-08) no such luck.

# etc
Miscellaneous code to make golang a little kinder to the programmer
