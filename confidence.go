package getrational

import "math"
import "fmt"

func evaluateConfidence(correct int, questions int, statedConfidence float64) string {
	var message string

	if correct < 0 || correct > questions {
		return "Assertion failed in evaluateConfidence: variable questions out of bounds!"
	}

	// Two-sided binomial test
	pLeft := binomCDF(float64(correct), float64(questions), statedConfidence) * 2
	pRight := (1 - binomCDF(float64(correct-1), float64(questions), statedConfidence)) * 2
	factor := (float64(correct) / float64(questions)) / statedConfidence

	if pLeft <= 0.05 {
		message = fmt.Sprintf("You're significantly overconfident (p < %g), with an estimated factor of %.1fX",
			roundP(pLeft), 1/factor)
	} else if pRight <= 0.05 {
		message = fmt.Sprintf("You're significantly underconfident (p < %g), with an estimated factor of %.1fX",
			roundP(pRight), factor)
	} else {
		message = fmt.Sprintf("You seem to be calibrated very well. Your quota of correct answers is statistically in line with your stated confidence level of %g. Congratulations, you're a shining example of rationality!", statedConfidence)
	}

	return message
}

func roundP(p float64) float64 {
	var cutoffs = [7]float64{1, 0.05, 0.01, 0.001, 1e-4, 1e-5, 1e-6}

	m := 0
	for m < len(cutoffs) {
		if p > cutoffs[m] {
			break
		}
		m = m + 1
	}

	return cutoffs[m-1]
}

func binomCDF(k float64, n float64, p float64) float64 {
	return 1.0 - betai(k+1, n-k, p)
}

func gammln(xx float64) float64 {
	cof := [6]float64{76.18009172947146, -86.50532032941677, 24.01409824083091, -1.231739572450155, 0.1208650973866179e-2, -0.5395239384953e-5}

	x := xx
	y := xx
	tmp := x + 5.5
	tmp = tmp - (x+0.5)*math.Log(tmp)
	ser := 1.000000000190015
	j := 0
	for j <= 5 {
		y = y + 1
		ser = ser + cof[j]/y
		j = j + 1
	}

	return -tmp + math.Log(2.5066282746310005*ser/x)
}

func betacf(a float64, b float64, x float64) float64 {
	var aa, c, d, del, h, qab, qam, qap float64
	const FPMIN float64 = 1.0e-30

	qab = a + b
	qap = a + 1.0
	qam = a - 1.0
	c = 1.0
	d = 1.0 - qab*x/qap

	if math.Abs(d) < FPMIN {
		d = FPMIN
	}
	d = 1.0 / d
	h = d

	m := 1
	for m <= 100 {
		m1 := float64(m)
		m2 := 2 * m1
		aa = m1 * (b - m1) * x / ((qam + m2) * (a + m2))
		d = 1.0 + aa*d
		if math.Abs(d) < FPMIN {
			d = FPMIN
		}
		c = 1.0 + aa/c
		if math.Abs(d) < FPMIN {
			d = FPMIN
		}
		d = 1.0 / d
		h *= d * c
		aa = -(a + m1) * (qab + m1) * x / ((a + m2) * (qap + m2))
		d = 1.0 + aa*d
		if math.Abs(d) < FPMIN {
			d = FPMIN
		}
		c = 1.0 + aa/c
		if math.Abs(d) < FPMIN {
			d = FPMIN
		}
		d = 1.0 / d
		del = d * c
		h *= del
		if math.Abs(del-1.0) < 3.0e-7 {
			break
		}
		m = m + 1
	}

	if m > 100 {
		return math.NaN()
	}
	return h
}

func betai(a float64, b float64, x float64) float64 {
	var bt float64

	if x < 0.0 || x > 1.0 {
		return math.NaN()
	}

	if x == 0.0 || x == 1.0 {
		bt = 0.0
	} else {
		bt = math.Exp(gammln(a+b) - gammln(a) - gammln(b) + a*math.Log(x) + b*math.Log(1.0-x))
	}

	if x < (a+1.0)/(a+b+2.0) {
		return bt * betacf(a, b, x) / a
	}
	return 1.0 - bt*betacf(b, a, 1.0-x)/b
}
