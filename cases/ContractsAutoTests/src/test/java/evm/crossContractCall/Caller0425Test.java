package evm.crossContractCall;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.Callee0425;
import network.platon.contracts.Caller0425;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title 0.4.25跨合约调用的调用者
 *  说明：CALL修改的是被调用者的状态变量，使用的是上一个调用者地址
 *       DELEGATECALL会一直使用原始调用者的地址，而CALLCODE不会。两者都是修改被调用者的状态
 * @description:
 * @author: hudenian
 * @create: 2019/12/30
 */
public class Caller0425Test extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "Caller0425Test-0.4.25跨合约调用CALL", sourcePrefix = "evm")
    public void caller0425CallTest() {
        try {
            //调用者合约地址
            Caller0425 caller0425 = Caller0425.load("0xcf9ed4c1670e53d6f457fc7b25f77a2d83015928",web3j, transactionManager, provider);
            String callerContractAddress = caller0425.getContractAddress();


            //被调用者合约地址
            Callee0425 callee0425 = Callee0425.load("0x0c3037b685421017981f82b21e1fb1e47756d106",web3j, transactionManager, provider);
            String calleeContractAddress = callee0425.getContractAddress();

            //查询调用者x值
            String callerX = caller0425.getCallerX().send().toString();
            collector.logStepPass("Caller0425 合约中X的值为："+callerX);

            //查询被调用者x值
            String calleeX = callee0425.getCalleeX().send().toString();
            collector.logStepPass("Callee0425 合约中X的值为："+calleeX);


            TransactionReceipt tx2 = caller0425.inc_call(calleeContractAddress).send();
            collector.logStepPass("执行跨合约调用后，hash:" + tx2.getTransactionHash());

            //查询调用者x值
            String callerAfterX = caller0425.getCallerX().send().toString();
            collector.logStepPass("跨合约调用后，Caller0425 合约中X的值为："+callerAfterX);

            //查询被调用者x值
            String calleeAfterX = callee0425.getCalleeX().send().toString();
            collector.logStepPass("跨合约调用后，Callee0425 合约中X的值为："+calleeAfterX);


        } catch (Exception e) {
            collector.logStepFail("Caller0425Test caller0425CallTest  process fail.", e.toString());
            e.printStackTrace();
        }
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test1.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "Caller0425Test-0.4.25跨合约调用CALLCODE", sourcePrefix = "evm")
    public void caller0425CallCodeTest() {
        try {
            //调用者合约地址
            Caller0425 caller0425 = Caller0425.load("0xc81e790282897cf5e1aeff2e3af25b2d41a4601b",web3j, transactionManager, provider);
            String callerContractAddress = caller0425.getContractAddress();


            //被调用者合约地址
            Callee0425 callee0425 = Callee0425.load("0xbf8a3a4a937efe73355bb997b7d1b667e585c7be",web3j, transactionManager, provider);
            String calleeContractAddress = callee0425.getContractAddress();

            //查询调用者x值
            String callerX = caller0425.getCallerX().send().toString();
            collector.logStepPass("Caller0425 合约中X的值为："+callerX);

            //查询被调用者x值
            String calleeX = callee0425.getCalleeX().send().toString();
            collector.logStepPass("Callee0425 合约中X的值为："+calleeX);


            TransactionReceipt tx2 = caller0425.inc_callcode(calleeContractAddress).send();
            collector.logStepPass("执行跨合约调用后，hash:" + tx2.getTransactionHash());

            //查询调用者x值
            String callerAfterX = caller0425.getCallerX().send().toString();
            collector.logStepPass("跨合约调用后，Caller0425 合约中X的值为："+callerAfterX);

            //查询被调用者x值
            String calleeAfterX = callee0425.getCalleeX().send().toString();
            collector.logStepPass("跨合约调用后，Callee0425 合约中X的值为："+calleeAfterX);


        } catch (Exception e) {
            collector.logStepFail("Caller0425Test caller0425CallCodeTest  process fail.", e.toString());
            e.printStackTrace();
        }
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test2.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "Caller0425Test-0.4.25跨合约调用DELEGATECALL", sourcePrefix = "evm")
    public void caller0425DelegateCallTest() {
        try {
            //调用者合约地址
            Caller0425 caller0425 = Caller0425.load("0x29a1fc4ea037f96dbff19df00ad15d86be1d268d",web3j, transactionManager, provider);
            String callerContractAddress = caller0425.getContractAddress();


            //被调用者合约地址
            Callee0425 callee0425 = Callee0425.load("0xe4389a17563c5c02a0b988e169d769bd529f7996",web3j, transactionManager, provider);
            String calleeContractAddress = callee0425.getContractAddress();

            //查询调用者x值
            String callerX = caller0425.getCallerX().send().toString();
            collector.logStepPass("Caller0425 合约中X的值为："+callerX);

            //查询被调用者x值
            String calleeX = callee0425.getCalleeX().send().toString();
            collector.logStepPass("Callee0425 合约中X的值为："+calleeX);


            TransactionReceipt tx2 = caller0425.inc_delegatecall(calleeContractAddress).send();
            collector.logStepPass("执行跨合约调用后，hash:" + tx2.getTransactionHash());

            //查询调用者x值
            String callerAfterX = caller0425.getCallerX().send().toString();
            collector.logStepPass("跨合约调用后，Caller0425 合约中X的值为："+callerAfterX);

            //查询被调用者x值
            String calleeAfterX = callee0425.getCalleeX().send().toString();
            collector.logStepPass("跨合约调用后，Callee0425 合约中X的值为："+calleeAfterX);


        } catch (Exception e) {
            collector.logStepFail("Caller0425Test caller0425DelegateCallTest  process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
