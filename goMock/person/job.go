package person

type Job interface {
  Get(id int64) error
}
