# Scripts

Folder containing some scripts to test or/and demo the functionality. Before executing any script you must have the `furyad` installed.

1. **init.sh** create the chain with the initial values
2. **start.sh** start the chain
3. **gov.sh** submit the gov.json governance proposal, votes on favor and query the created furya
4. **delegate.sh** delegate to the previously create furya and query the modified furya
5. **rewards.sh** claim rewards and query information about the evidences of the process
6. **undelegate.sh** undelegante the tokens from the furya and query the evidences
7. **gov-del.sh** with the file gov-delete.json deletes the furya created in third step.

> This scripts must be executed in the specified order since they have dependencies on each other.
