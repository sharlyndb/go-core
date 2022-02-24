/**
 * @Time: 2022/2/24 10:42
 * @Author: yt.yin
 */

package srun

type server interface {
	ListenAndServe() error
}

