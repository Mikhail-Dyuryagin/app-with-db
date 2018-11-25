pg_dump -h localhost -O -F t -c -U devops  api | gzip -c > myDB-filedump.gz   -   export
pg_dumpall | gzip -c > allDB-filedump.gz   -   import

