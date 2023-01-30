# Import the DBMS
source Simplest_db.sh

# Configuration
large_string_length=100
#medium_string_length=50
#small_string_length=20

load_upper_bound=100000
load_lower_bound=10000
load_increment=10000

test_get_number_of_repetitions=10000

stats_filename='Simplest_db_eval_stats'
read_stats_filename='Simplest_db_eval_read_db'
# Pass Length of the string as parameter
gen_string() {
  openssl rand -hex "$1"
}

# Testing
before() {
  : > "$stats_filename"
  record_stats 'TestName' 'InputSize' 'Runtime(seconds)' 'Latency' 'Throughput'
}

before_test() {
  function_name=$1
  db_truncate
  if [[ "$function_name" == 'test_get' ]]
  then
    input_size=$2
    echo "$input_size"
    # shellcheck disable=SC2154
    cat "$read_stats_filename" | head -n "$input_size" > "$db_name"
  fi
}

# With variable Input sizes call Get 10000 times
# Parameter: Size of DB
test_get() {
  input_size=$1
  i=0
  while [[ $i -lt "$test_get_number_of_repetitions" ]]
  do
    id=$(jot -r 1 0 "$input_size")
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
    throughput=$(echo 4k $work $runtime /p | dc)
    latency=$(echo 4k $runtime $work /p | dc)
    record_stats "$test_function" "$load_counter" "$runtime" "$latency" "$throughput"
    load_counter=$((load_counter + load_increment))
  done
#  after
}

record_stats() {
  echo "$1, $2, $3, $4, $5" >> "$stats_filename"
}
