
(define conso
    (lambda (a d p)
        (== (cons a d) p)))
(define pairo
    (lambda (p)
        (fresh (a d)
            (conso a d p))))

(test-check "pairo"
    (run* (q) (pairo q))
    '((_.0 . _.1)))

(test-check "let unify #1"
    (run* (q)
        (let ((a (== #t q))
                (b (fresh (x)
                (== x q)
                (== #f x)))
                (c (conde
                ((== #t q) succeed)
                (else (== #f q)))))
            b))
    '(#f))

(test-check "let unify #2"
    (run* (r)
        (fresh (v w)
            (== (let ((x v) (y w))
                    (list x y))
                r)))
    '((_.0 _.1)))

(test-check "salad"
    (run* (r)
        (fresh (x y)
            (== (cons x (cons y 'salad)) r)))
    '((_.0 _.1 . salad)))

(define teacupo
    (lambda (x)
        (conde
            ((== 'tea x) succeed)
            ((== 'cup x) succeed)
            (else fail))))

(test-check "teacupo #1"
    (run* (x) (teacupo x))
    '(tea cup))

(test-check "teacupo #2"
    (run* (r)
        (fresh (x y)
            (conde
                ((teacupo x) (== #t y) succeed)
                ((== #f x) (== #t y))
                (else fail))
            (== (cons x (cons y '())) r)))
    '((tea #t) (cup #t) (#f #t)))

(define appendo
    (lambda (l s out)
        (conde
            ((== '() l) (== s out))
            ((fresh (a d)
                (== (cons a d) l)
                (fresh (res)
                    (== (cons a res) out)
                    (appendo d s res)))))))

;; Byrd (2009), p. 19
(test-check "appendo"
    (run 6 (q)
        (fresh (l s)
            (appendo l s '(a b c d e))
            (== (list l s) q)))
    '((() (a b c d e)) ((a) (b c d e)) ((a b) (c d e)) ((a b c) (d e)) ((a b c d) (e)) ((a b c d e) ())))

(test-check "triple fresh"
    (run* (q)
        (fresh (x y z)
            (== q `(,x ,y ,z))
            (fresh (z)
                (conde
                    ((== z 'a))
                    ((== z 'b)))
                (== z x))
            (fresh (w)
                (conde
                    ((== w 'A))
                    ((== w 'B)))
                (== w y))
             (fresh (v)
                (conde
                    ((== v 1))
                    ((== v 2)))
                (== v z))))
    '((a A 1) (a A 2) (a B 1) (a B 2) (b A 1) (b A 2) (b B 1) (b B 2)))
