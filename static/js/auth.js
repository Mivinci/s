!function () {
  'use-strict'
  
  class Timer {
    constructor(el, wait) {
      this.el = el;
      this.timer = null;
      this.count = wait
    }

    start(timeout) {
      this.timer = setInterval(() => {
        if (this.count <= 0) {
          this.stop();
          this.count = this.wait;
          this.reset();
          return;
        }
        this.el.innerText = `${--this.count} 秒后可重发`;
      }, timeout);
    }

    reset() {
      this.el.innerText = "重新获取验证码";
    }

    stop() {
      clearInterval(this.timer);
      this.timer = null;
      this.reset();
    }

    get running() {
      return this.timer !== null;
    }
  }


  const codeBtn = document.getElementById("code-btn");
  const toInput = document.getElementById("tel-input");

  let timer = new Timer(codeBtn, 60);

  codeBtn.addEventListener("click", () => {
    const to = toInput.value;
    if (timer.running || !to) {
      return;
    }
    timer.start(1000);
    const uri = `/transport?to=${to}`;
    fetch(uri)
      .then(res => {
        if (res.status != 200) {
          timer.stop();
          throw new Error(res.statusText);
        }
        return res.json()
      })
      .then(({error}) => {
        console.log(error)
        if (error) {
          timer.stop();
          throw new Error(error);
        }
      })
      .catch(err => {
        alert(err.message);
      })
  });

}();