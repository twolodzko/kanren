;; Examples from Byrd & Friedman (2006)

(define add1
   (lambda (x)
      (list 's x)))

(define int->peano (lambda (x)
   (define f
      (lambda (n acc)
         (cond
            ((< n 0) #f)
            ((= n 0) acc)
            (else (f (- n 1) (add1 acc))))))
   (f x 'z)))

(define peano->int
   (lambda (x)
      (cond
         ((= x 'z) 0)
         (else (+ 1 (peano->int (car (cdr x))))))))

(define pluso
   (lambda (n m sum)
      (conde
         ((== 'z n) (== m sum))
         ((fresh (x y)
            (== (add1 x) n)
            (== (add1 y) sum)
            (pluso x m y))))))

(define minuso
   (lambda (n m k)
      (pluso m k n)))

(define eveno
   (lambda (n)
      (conde
         ((== 'z n))
         ((fresh (m)
            (== (add1 (add1 m)) n)
            (eveno m))))))

(define positiveo
   (lambda (n)
      (fresh (m)
         (== (add1 m) n))))

(define plus*o
   (lambda (in* out)
      (conde
         ((== '() in*) (== 'z out))
         ((fresh (a d res)
            (== (cons a d) in*)
            (pluso a res out)
            (plus*o d res))))))

;; ============= helpers =============

(load "examples/stdlib.scm")

(define peano?
   (lambda (x)
      (or
         (= x 'z)
         (and
            (list? x)
            (= (car x) 's)))))

(define map-peano->int
   (lambda (x)
      (cond
         ((peano? (car x)) (map peano->int x))
         (else (map map-peano->int x)))))

;; ============= tests =============

(let ((x 0))
   (test-check "peano conversion for 0"
      (peano->int (int->peano x))
      x))
(let ((x 79))
   (test-check "peano conversion for 79"
      (peano->int (int->peano x))
      x))

(let ((x 0))
   (test-check "peano? for 0"
      (peano? (int->peano x))
      #t))
(let ((x 1))
   (test-check "peano? for 1"
      (peano? (int->peano x))
      #t))
(let ((x 2))
   (test-check "peano? for 2"
      (peano? (int->peano x))
      #t))
(let ((x 173))
   (test-check "peano? for 173"
      (peano? (int->peano x))
      #t))

(test-check "pluso"
   (map-peano->int
      (run* (q)
         (fresh (n m)
            (pluso n m (int->peano 6))
            (== (list n m) q))))
   '((0 6) (1 5) (2 4) (3 3) (4 2) (5 1) (6 0)))

(test-check "minuso"
   (map-peano->int
      (run 10 (q)
         (fresh (n m)
            (minuso n m (int->peano 6))
            (== (list n m) q))))
   '((6 0) (7 1) (8 2) (9 3) (10 4) (11 5) (12 6) (13 7) (14 8) (15 9)))

(test-check "eveno"
   (map-peano->int
      (run 4 (q) (eveno q)))
   '(0 2 4 6))

(test-check "plus*o #1"
   (run* (q)
      (plus*o (list (int->peano 3) (int->peano 4) (int->peano 2)) q))
   (list (int->peano 9)))

(test-check "plus*o #2"
   (run* (q) (plus*o (cons (int->peano 5) q) (int->peano 3)))
   '())
