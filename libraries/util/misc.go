package util

import (
	"errors"
	"math/rand"
	"time"
)

/**
 * Difference by default returns only elems in l1 but not in l2.
 * When bidirectional is true, all elements that are in this set but not the others will be returned.
 */
func Difference(l1, l2 []int, bidirectional bool) []int {
	var (
		diff []int
		ok   bool
	)

	for _, v := range l1 {
		ok = false
		for _, j := range l2 {
			if v == j {
				ok = true
			}
		}

		if !ok {
			diff = append(diff, v)
		}
	}

	if bidirectional {
		for _, v := range l2 {
			ok = false
			for _, j := range l1 {
				if v == j {
					ok = true
				}
			}

			if !ok {
				diff = append(diff, v)
			}
		}
	}

	return diff
}

/**
* w: non-negative weight
* l: related items
* l[i]'s weight is w[i]
*
* return the chosen item
* RandomByWeight almost simplify to MultiRandomByWeight(l, 1)
 */
func RandomByWeight(w, l []int) (int, error) {
	if len(w) != len(l) {
		return 0, errors.New("weight list's length must be equal to item list's length")
	}

	total := 0
	for _, v := range w {
		if v <= 0 {
			return 0, errors.New("weight should be non-negative")
		}
		total += v
	}
	if total <= 0 {
		return 0, errors.New("total weight should be non-negative")
	}

	var (
		chosen int
		p      int
		c      int
	)

	rand.Seed(time.Now().Unix())
	p = rand.Intn(total)
	c = 0
	for idx, w := range w {
		if p >= c && p < c+w {
			chosen = idx
			break
		}
		c = c + w
	}

	return l[chosen], nil
}

/**
* w: non-negative weight
* l: related items
* l[i]'s weight is w[i]
*
* return the chosen item
 */
func MultiRandomByWeight(w, l []int, max_chosen_cnt int) ([]int, error) {
	if len(w) != len(l) {
		return nil, errors.New("weight list's length must be equal to item list's length")
	}

	total := 0
	for _, v := range w {
		if v <= 0 {
			return nil, errors.New("weight should be non-negative")
		}
		total += v
	}
	if total <= 0 {
		return nil, errors.New("total weight should be non-negative")
	}

	if len(l) <= max_chosen_cnt {
		return l, nil
	}

	var (
		chosen       = make([]int, max_chosen_cnt)
		chosen_flags = make([]bool, len(l))

		p int
		c int
	)

	rand.Seed(time.Now().Unix())
	for round := 0; round < max_chosen_cnt; round += 1 {
		p = rand.Intn(total)
		c = 0
		for idx := 0; idx < len(w); idx++ {
			if !chosen_flags[idx] {
				if p >= c && p < c+w[idx] {
					chosen_flags[idx] = true
					chosen[round] = l[idx]
					total -= w[idx]
					break
				}
				c = c + w[idx]
			}
		}
	}

	return chosen, nil
}

/**
 * Cn,m returns Combination n,m.
 * n >= m
 */
func cnm(n, m int) int {
	return anm(n, m) / factorial(m)
}

/**
 * An,m returns Arrangement n,m.
 * n >= m
 */
func anm(n, m int) int {
	r := 1
	m = n - m
	for n > m {
		r *= n
		n -= 1
	}
	return r
}

/**
 * factorial n returns n!.
 * 0!=1
 */
func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}
