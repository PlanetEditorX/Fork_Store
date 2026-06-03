(function () {
    window.Game = function () {
        //设置行和列
        this.row = 20
        this.col = 12
        //初始化
        this.init()
        //实例方块
        this.block = new Block()
        //实例下一个方块
        this.nextBlock = new Block()
        //实例地图
        this.map = new Map()
        //启动定时器
        this.start()
        //键盘事件监听
        this.bindEvent()
        //分数
        this.score = 0
        //速度
        this.during = 30
        //帧
        this.f = 0
    }
    //初始化布局
    Game.prototype.init = function () {
        //初始化大表格
        var $table = $("<table></table>")
        $table.addClass('tab1')
        //渲染表格
        for (var i = 0; i < this.row; i++) {
            //创建tr
            var $tr = $("<tr></tr>")
            for (var j = 0; j < this.col; j++) {
                //创建td
                var $td = $("<td></td>")
                //将每个单元格追加到行
                $td.appendTo($tr)
            }
            //将行追加到表格
            $tr.appendTo($table)
        }
        //初始化预览窗口
        var $table2 = $("<table></table>")
        $table2.addClass('tab2')
        for (let i = 0; i < 4; i++) {
            var $tr2 = $("<tr></tr>")
            for (let j = 0; j < 4; j++) {
                var $td2 = $("<td></td>")
                $td2.appendTo($tr2)
            }
            $tr2.appendTo($table2)
        }
        //将表格追加到body
        $($table).appendTo("body")//主表格
        $($table2).appendTo("body")//预览表格

    }

    /**
     * 方块渲染
     * @param {int} row 行
     * @param {int} col 列
     * @param {int} num 状态
     */
    Game.prototype.setColor = function (row, col, num) {
        // 给对应的有颜色方块添加类名
        $('.tab1').find('tr').eq(row).children("td").eq(col).addClass("c" + num)
    }
    Game.prototype.setNextColor = function (row, col, num) {
        for (let i = 0; i < 4; i++) {
            for (let j = 0; j < 4; j++) {
                if (this.nextBlock.code[i][j] != 0) {
                    $('.tab2').find('tr').eq(i).children("td").eq(j).addClass("c" + this.nextBlock.code[i][j])
                }
            }
        }
    }
    //清屏功能
    Game.prototype.clear = function () {
        for (let i = 0; i < this.row; i++) {
            for (let j = 0; j < this.col; j++) {
                //移除所有的class
                $(".tab1").find('tr').eq(i).children("td").eq(j).removeClass()
            }
        }
        for (let i = 0; i < 4; i++) {
            for (let j = 0; j < 4; j++) {
                //移除所有的class
                $(".tab2").find('tr').eq(i).children("td").eq(j).removeClass()
            }
        }

    }
    //定时器
    Game.prototype.start = function () {
        var self = this
        //帧编号
        this.f = 0
        this.timer = setInterval(function () {
            self.f++
            //渲染帧编号
            document.getElementById("f").innerHTML = `帧编号: ${self.f}`
            //清屏
            self.clear()
            //渲染方块
            self.block.render()
            //渲染预览方块
            self.setNextColor()
            //渲染地图
            self.map.render(self)
            //下落
            self.f % self.during === 0 && self.block.checkDown()
        }, 20)
    }
    //事件监听
    Game.prototype.bindEvent = function () {
        //备份
        let self = this
        document.onkeydown = function (event) {
            if (event.key === 'ArrowLeft') {
                //判断是否能向左走
                self.block.checkLeft()
            } else if (event.key === 'ArrowRight') {
                self.block.checkRight()
            } else if (event.key === 'ArrowUp') {
                //用来切换方块方向
                self.block.checkRot()
            } else if (event.key === 'ArrowDown') {
                self.block.checkBlockEnd()
            }
        }
    }
})()