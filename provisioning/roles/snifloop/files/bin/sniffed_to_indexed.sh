cat sniffed_hashes.txt | sort | uniq -c | sort -n | awk '{ if ($1 > 2) {print $2} }' >> indexed_hashes.txt

