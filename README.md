# Abe
Abe, a filesystem changes monitoring tool, much like Aide or Tripwire with the 
difference that it stores hashes in a Postgresql database and reports changes 
through email and syslog.
It uses filesystem watches for detecting changes.
Some more work has to be done on the loadable kernel modules checks and
several network port checks. (definitely not ellegant)

BTW it's far from finished, but it works....
