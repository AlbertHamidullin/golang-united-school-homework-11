package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

type safeRes struct {
	m    sync.Mutex
	data []user
}

func newSafeRes() safeRes {
	return safeRes{data: make([]user, 0)}
}
func (s *safeRes) Append(u user) {
	s.m.Lock()
	defer s.m.Unlock()
	s.data = append(s.data, u)
}
func (s *safeRes) Len() int64 {
	s.m.Lock()
	defer s.m.Unlock()
	return int64(len(s.data))
}
func (s *safeRes) GetData() []user {
	s.m.Lock()
	defer s.m.Unlock()
	dest := make([]user, len(s.data))
	copy(dest, s.data)
	return dest
}

func getBatch(n int64, pool int64) []user {
	var sem = make(chan struct{}, pool)
	var sRes safeRes = newSafeRes()
	var wg sync.WaitGroup

	for i := int64(0); i < n; i++ {
		sem <- struct{}{}
		wg.Add(1)
		go func(sem chan struct{}, wg *sync.WaitGroup, sRes *safeRes, i int64) {
			defer wg.Done()
			sRes.Append(getOne(i))
			<-sem
		}(sem, &wg, &sRes, i)
	}

	close(sem)

	wg.Wait()

	return sRes.GetData()
}

// func getBatch6(n int64, pool int64) []user {
// 	var sem = make(chan int64, pool)
// 	var sRes safeRes = newSafeRes()
// 	var wg sync.WaitGroup

// 	for i := int64(0); i < n; i++ {
// 		sem <- i
// 		wg.Add(1)
// 		go func(i int64) {
// 			defer wg.Done()
// 			sRes.Append(getOne(i))
// 			<-sem
// 		}(i)
// 	}

// 	close(sem)

// 	wg.Wait()

// 	return sRes.GetData()
// }

// func getBatch5(n int64, pool int64) (res []user) {
// 	var sem = make(chan int64, pool)
// 	var wg sync.WaitGroup

// 	res = make([]user, n)

// 	for i := int64(0); i < n; i++ {
// 		sem <- i
// 		wg.Add(1)
// 		go func(i int64) {
// 			defer wg.Done()
// 			res[i] = getOne(i)
// 			<-sem
// 		}(i)
// 	}

// 	close(sem)

// 	wg.Wait()

// 	return res
// }

// func getBatch4(n int64, pool int64) (res []user) {
// 	var sem = make(chan int64, pool)

// 	res = make([]user, n)

// 	for i := int64(0); i < n; i++ {
// 		sem <- i
// 		go func(i int64) {
// 			res[i] = getOne(i)
// 			<-sem
// 		}(i)
// 	}

// 	for len(sem) > 0 {
// 		time.Sleep(time.Millisecond)
// 	}

// 	return res
// }

// func getBatch3(n int64, pool int64) (res []user) {
// 	var sem = make(chan int64, pool)
// 	var count = make(chan int64, n)

// 	res = make([]user, n)

// 	for i := int64(0); i < n; i++ {
// 		sem <- i
// 		go func(i int64) {
// 			res[i] = getOne(i)
// 			<-sem
// 			count <- 0
// 		}(i)
// 	}

// 	for i := int64(0); i < n; i++ {
// 		<-count
// 	}

// 	return res
// }

// func getBatch2(n int64, pool int64) (res []user) {
// 	res = make([]user, n)

// 	for i := int64(0); i < n; i++ {
// 		res[i] = getOne(i)
// 	}

// 	return res
// }
