(define reverse
   (lambda (l)
      (define impl
         (lambda (l acc)
            (cond
               ((null? l) acc)
               (else
                  (impl
                     (cdr l)
                     (cons (car l) acc))))))
      (impl l '())))

(define map
   (lambda (f x)
      (define impl
         (lambda (f x acc)
            (cond
               ((null? x) (reverse acc))
               (else
                  (impl
                     f
                     (cdr x)
                     (cons (f (car x)) acc))))))
      (impl f x '())))

(define list?
   (lambda (x)
      (cond
         ((null? x) #t)
         ((pair? x) (list? (cdr x)))
         (else #f))))

(define eq?
   (lambda (a b)
      (= a b)))

(define length
   (lambda (l)
      (cond
         ((null? l) 0)
         (else (+ 1 (length (cdr l)))))))

(define zero? (lambda (x) (= x 0)))
(define even? (lambda (x) (zero? (% x 2))))
(define odd? (lambda (x) (not (even? x))))
(define quotient (lambda (x y) (/ x y)))

(test-check "list? #1" (list? '()) #t)
(test-check "list? #2" (list? '(1 2 3)) #t)
(test-check "list? #3" (not (list? '(1 . 2))) #t)
(test-check "list? #4" (not (list? #t)) #t)
