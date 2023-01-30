db_name='Simplest_db'
tombstone='`' #Placeholder for deleting a key
# Meaning: Function takes key and value, then stores them in a comma-separated style in the Simplest_db file
db_set() {
  echo "$1, $2" >> $db_name
}

# Meaning:
# Function takes the key value
# Grep - Finds all lines that start with the key value
# Sed - Remove the key part from all lines
# Tail - Take the last occurrence (i.e last line) of the output

# Then, we have to check that:
# A) Value has been set
# B) Value is not "tombstone", that is value is the deletion of record, we do not show it
db_get() {
  out=$(grep "^$1," $db_name | sed "s/^$1, //" | tail -n 1)
  if [ -n "$out" ] && [ "$out" != "$tombstone" ]; then
    echo "$out"
  fi
}

# Meaning:
# This function clears the contents of the database
# : is a no-op in Bash
# this essentially opens the file for writing (which clears the file) and then closes it immediately
db_truncate() {
  : > $db_name
}

# Meaning:
# Deleting an item requires just setting am item value to a pre-defined "Tombstone", mentioned above
db_delete() {
  db_set "$1" "$tombstone"
}