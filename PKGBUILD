pkgname=sapfun
pkgver=0
pkgrel=1
pkgdesc='Utility that takes control over your video card coolers to keep it cool and steady.'
arch=('x86_64')
url='https://git.vimusov.space/me/sapfun'
license=('GPL')
makedepends=('go' 'make')
source=("${pkgname}.go" makefile service)
md5sums=('SKIP' 'SKIP' 'SKIP')

pkgver()
{
    printf "r%s.%s" "$(git rev-list --count HEAD)" "$(git rev-parse --short HEAD)"
}

build()
{
    make -C "$srcdir"
}

package()
{
    make -C "$srcdir" DESTDIR="$pkgdir" install
    install -D --mode=0644 "$srcdir"/service "$pkgdir"/usr/lib/systemd/system/"${pkgname}.service"
}
