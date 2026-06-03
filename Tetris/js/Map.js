(function () {
    //地图的矩阵
    window.Map = function () {
        // console.log(mapGame);
        this.mapCode = [
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9]//数据墙
        ]

    }
    Map.prototype.render = function (mapGame) {
        //渲染地图
        for (let i = 0; i < mapGame.row; i++) {
            for (let j = 0; j < mapGame.col; j++) {
                if (this.mapCode[i][j] != 0) {
                    game.setColor(i, j, this.mapCode[i][j])
                }
            }
        }
    }
    Map.prototype.checkRemove = function () {
        //判断当前的mapCode是否能够消行
        //消行规则: 当前的mapCode数组的每一项如果都不是0,则说明该消行了
        for (let i = 0; i < 20; i++) {
            //遍历地图数组
            if (this.mapCode[i].indexOf(0) == -1) {
                //删除这一行
                this.mapCode.splice(i, 1)
                //补一行
                this.mapCode.unshift([0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0])
                //分数增加
                game.score += 10
                if (game.during <= 30 && game.during >= 20) {
                    game.score += 10
                } else if (game.during < 20 && game.during >= 10) {
                    game.score += 20
                } else {
                    game.score += 30
                }
                //渲染分数
                document.getElementById("score").innerHTML = `分数: ${game.score}`
                if (game.score % 100 === 0) {
                    game.during -= 5
                    if (game.during <= 0) {
                        game.during = 1
                    }
                }
            }
        }
    }
})()