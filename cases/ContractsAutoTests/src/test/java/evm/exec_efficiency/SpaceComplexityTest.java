package evm.exec_efficiency;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.SpaceComplexity;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title SpaceComplexityTest
 * @description 空间复杂度场景测试
 * @author qcxiao
 * @updateTime 2019/12/28 14:39
 */
public class SpaceComplexityTest extends ContractPrepareTest {
    private BigInteger numberOfCalls;
    private String contractAddress;

    @Before
    public void before() {
        numberOfCalls = new BigInteger(driverService.param.get("numberOfCalls"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "exec_efficiency.SpaceComplexityTest-空间复杂度", sourcePrefix = "evm")
    public void test() {
        prepare();
        try {
            SpaceComplexity spaceComplexity = SpaceComplexity.load("0xd068caa72c9f1296d05035554a48e87b246b6827",web3j, transactionManager, provider);
            contractAddress = spaceComplexity.getContractAddress();

            TransactionReceipt transactionReceipt = SpaceComplexity.load(contractAddress, web3j, transactionManager, provider)
                    .testStorage(numberOfCalls).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String name = SpaceComplexity.load(contractAddress, web3j, transactionManager, provider).name().send();

            if (numberOfCalls.mod(BigInteger.valueOf(2)) == BigInteger.ZERO) {
                collector.assertEqual(name, "QCXIAO");
            } else {
                collector.assertEqual(name, "qcxiao");
            }
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}
