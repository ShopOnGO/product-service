package interfaces

type ProductChecker interface {
    ExistsByID(id uint) (bool, error)
}