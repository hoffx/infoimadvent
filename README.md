# infoimadvent
A dynamic website for students of all ages to have fun before Christmas

## Development

To get this running, go get this repository or git pull it:

`````
go get github.com/hoffx/infoimadvent
`````
or
`````
git clone https://github.com/hoffx/infoimadvent.git
`````

Afterwards go to the project root and run the makefile:
`````
make
`````


Now, you have to set up a mysql-database and a mail-server for this project. Of course you can also use an external one for both.


Then you have to create your configuration file:


Create a `config.ini`file and copy the `config.example.ini`'s content there. The comments should help you get that done. The only things you should have to change are the `server ip`, `port` and `address`, as well as the `db` and `mail` sections.

## Production

### Get it running

Download the latest release for your OS and architecture from the [releases](https://github.com/hoffx/infoimadvent/releases) page and unzip it.


Now, you have to set up a mysql-database and a mail-server for this project. Of course you can also use an external one for both.


The downloaded folder should contain some configuration file with the extension `.ini`. Open and edit it according to your needs. The most relevant sections are `server`, `db` and `mail`. Also don't forget to change the hash for your admin-password. You can find it in the `auth` section.

Now start the webserver by typing, depending on your server's operating system, ether (linux)

`````
./infoimadvent web
`````
or (windows)
`````
infoimadvent.exe web
`````

If you didn't mess up anything (e.g. with the ports) the service should now be available using the `ip` and `port` provided in the config-file.

### Fill it with content

Now, the service is running, but the whole calendar-section isn't ready for public yet. That's because you haven't uploaded any entries yet.


To do so, head over to `/login` and log in using your admin credentials defined in the config-file. It the login was successful you are in an admin-session now. So now you can access the hidden pages `/upload` and `/overview`.


The overview will print out a table of all days of the advent and all grades. In each cell is a string representing the state of the according document. While `+++` signalizes the document was uploaded is working correctly, `!/-` stands for a document where some error occurred. The `---` states out, that the according document hasn't been uploaded yet.


On the upload page you can upload a new document. The upload form has the following fields:

* `type` can ether be `Terms of Service`, `About` or `Quest`. The last one is needed for uploading a question for the calendar. This is the only field next to `markdown` which is always required. If you choose type `Quest` the following fields are required additionally: `min grade`, `max grade`, `day` and `solution`.
* `markdown` is used to select the Markdown (.md) file containing the documents content. You can find a cheatsheet for Markdown [here](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet).
* `assets zip` is used to select the documents assets packed as a zip (.zip) archive. If you don't specify one, a warning is going to show up after submitting, but the document is still going to work as intended in case it doesn't require any assets.
* `solution` can be `A`, `B`, `C` or `D`. You have to specify the correct answer to your question there if your document is of type `Quest`.
* `day` specifies the day of advent the question is going to be shown on.
* `max grade` and `min grade` are the boundaries for the grades the question will be shown on. For example, if the service's min-grade is 1 and the max-grade is 12, you can have questions, which are the same for all grades (by setting the form's `min grade` to 1 and the `max grade` to 12), or you can upload different questions for different grades (e.g. one from 1 to 4, one from 5 to 5 and one from 6 to 12).

**Just don't mess around with this form, because it doesn't check that much for nonsense input and you probably don't want to risk having to reset the server. Also, if you specify the same grade and day for a quest-document twice, the second upload will overwrite the first one without warning. So check the overview before uploading documents !**

### Maintenance

The server should reset itself at the specified `resetmonth` (config-file) and delete all questions and users (except admin) and calculate the scores correctly. However, if this fails (e.g. because a task was scheduled, but the server went down at that time as a result of a mistake) you can manually do it via the command line using the `./infoimadvent calc --day [day]` and `./infoimadvent reset --standard` commands. Type `./infoimadvent [command] --help` for further information.


**Good Luck !**