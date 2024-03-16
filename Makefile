maelstrom_path = ~/maelstrom/maelstrom

serve:
	$(maelstrom_path) serve

echo:
	go build -o ./bin ./cmd/echo

unique-ids:
	go build -o ./bin ./cmd/unique-ids

broadcast:
	go build -o ./bin ./cmd/broadcast

1: echo
	$(maelstrom_path) test -w echo --bin bin/echo --node-count 1 --time-limit 10

2: unique-ids
	$(maelstrom_path) test -w unique-ids --bin bin/unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

3a: broadcast
	$(maelstrom_path) test -w broadcast --bin bin/broadcast --node-count 1 --time-limit 20 --rate 10

3b: broadcast
	$(maelstrom_path) test -w broadcast --bin bin/broadcast --node-count 5 --time-limit 20 --rate 10

.PHONY: serve echo unique-ids broadcast 1 2 3a