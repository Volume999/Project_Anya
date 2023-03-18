# Import the DBMS
source BashDB.sh

# Documentation
# In this script we measure the performance of the System Using 2 metrics
# Throughput and Latency
# Writing: Increasing the number of records to write
# Reading; Increasing the size of the DB over which to perform reads
# Performance measuring done in style of Property-tests with generators and evaluators

# Running Evaluations:
# If you are planning to run Read test, it is better to setup a Read DB from which to copy data
# $ setup_read_db
# Then execute test function with the name of the test
# $ test test_get or test test_set_string
# Evaluation will appear in the stats file

# Configuration

# String size for insertion
large_string_length=100
#medium_string_length=50
#small_string_length=20

# For write: How many times to write to a disk
# For reads: On a DB of which size to execute reads
load_upper_bound=100000
load_lower_bound=10000
load_increment=10000

# For reads: How many times will a read function be executed
test_get_number_of_repetitions=10000

# File where the stats will be posted
stats_filename='Simplest_db_eval_stats.csv'

# Optimization: Backup DB from which a dataset will be formed for GET function (otherwise it is very slow)

read_stats_filename='Simplest_db_eval_read_db'

# Function generates a random string
# Parameter: Length of string
gen_string() {
  openssl rand -hex "$1"
}

# Function sets up a backup DB for read tests
setup_read_db() {
  i=0
  db_truncate
  while [[ "$i" -lt "$load_upper_bound" ]]
  do
    v=$(gen_string "$large_string_length")
    db_set "$i" "$v"
    ((i = i + 1))
  done
  # shellcheck disable=SC2154
  cat "$db_name" > "$read_stats_filename"
  db_truncate
}

# Prepare Statfile for testing
before() {
  : > "$stats_filename"
  record_stats 'TestName' 'InputSize' 'Runtime(seconds)' 'Latency' 'Throughput'
}

# Prepare environment for each iteration of a test
before_test() {
  function_name=$1
  db_truncate
  if [[ "$function_name" == 'test_get' ]]
  then
    input_size=$2
    echo "$input_size"
    # shellcheck disable=SC2154
    cat "$read_stats_filename" | head -n "$input_size" > "$db_name" # Copy the Input size number of records to the database
  fi
}

# With variable Input sizes call Get 10000 times
# Parameter: Size of DB
test_get() {
  input_size=$1
  i=0
  while [[ $i -lt "$test_get_number_of_repetitions" ]]
  do
    id=$(jot -r 1 0 "$input_size") # Generate random ID
    db_get "$id" &> /dev/null
    ((i = i + 1))
  done
}

# Sets a large string some number of times
# Parameter: Input size
test_set_string() {
  input_size=$1
  i=0
  while [[ $i -lt $input_size ]]
  do
    v=$(gen_string "$large_string_length")
    db_set "$i" "$v"
    ((i = i + 1))
  done
}

# Parameter: Test function name
test() {
  load_counter=$load_lower_bound
  test_function=$1
#  before
  while [ $load_counter -le $load_upper_bound ]
  do
#    stats=50
    before_test "$test_function" "$load_counter"
    start=$( (date +%s) )
    "$test_function" "$load_counter"
    end=$( (date +%s) )
    runtime=$((end - start))
    work=$load_counter
    if [[ "$test_function" == 'test_get' ]]
    then
      work=$test_get_number_of_repetitions
    fi
    throughput=$(echo 4k $work $runtime /p | dc) # Using 4 precision point, divide work by runtime
    latency=$(echo 4k $runtime $work /p | dc) # Using 4 precision point, divide runtime by work
    record_stats "$test_function" "$load_counter" "$runtime" "$latency" "$throughput"
    load_counter=$((load_counter + load_increment))
  done
#  after
}

# Parameters: Test name, Input size, Runtime, Latency, Throughput
record_stats() {
  echo "$1, $2, $3, $4, $5" >> "$stats_filename"
}
