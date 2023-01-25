# Meaning: Function takes key and value, then stores them in a comma-separated style in the Simplest_db file
db_set() {
  echo "$1, $2" >> Simplest_db
}

# Meaning:
# Function takes the key value
# Grep - Finds all lines that start with the key value
# Sed - Remove the key part from all lines
# Tail - Take the last occurrence (i.e last line) of the output
db_get() {
  grep "^$1," Simplest_db | sed "s/^$1, //" | tail -n 1
}