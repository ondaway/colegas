Person command
==============

This is a CLI (command line interface) for *person service*.

## Usage

    $ ./person [options] <personid>


Examples:

    # Search the person with id=jomoespe 
    $ ./person jomoespe

    # Search the person with id=jomoespe and then get the surname (the third field) 
    $ ./person juergas | cut -d$'\t' -f3
