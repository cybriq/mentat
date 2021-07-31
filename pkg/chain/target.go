package chain

// A Target is calculated by the difficulty adjustment algorithm and defines
// the minimum timestamp and the center point of the timestamp displacement
// adjustment
type Target struct {
	min, target int64
}
